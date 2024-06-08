import * as vscode from "vscode";
import ollama from "ollama";
import { AsyncLock } from "./asynclock";

import { GenerateRequest } from "ollama";
import type { GitExtension } from "./git";
import { activateDecorators } from "./decorator";
import {
  logAcceptedSuggestion,
  addAcceptsFromBlaimFile,
  AcceptLogLine,
} from "./acceptlog";

const modelName = "codegemma";
const acceptSuggestCommand = "blaim.acceptSuggestion";

function getPromptData(
  document: vscode.TextDocument,
  position: vscode.Position,
) {
  // Load document text
  const text = document.getText();
  const offset = document.offsetAt(position);
  const prefix = text.slice(0, offset);
  const suffix: string = text.slice(offset);
  return { prefix, suffix };
}

// The format for the prompt varies from model to model.
// TODO: clean this up so it uses a config type that abstracts out
// the model name, prompt format, temperature, number of tokens etc.
function formatPrompt(model: string, prefix: string, suffix: string) {
  if (model === "codellama") {
    return {
      prompt: `<PRE> ${prefix} <SUF> ${suffix} <MID>`,
      stop: [`<END>`, `<EOD>`, `<EOT>`],
    };
  } else if (model === "codegemma") {
    return {
      prompt: `<|fim_prefix|>${prefix}<|fim_suffix|>${suffix}<|fim_middle|>`,
      stop: ["<|file_separator|>"],
    };
  }
  return {
    prompt: prefix,
  };
}

const lock = new AsyncLock();

export async function loadBlaimFile(context: vscode.ExtensionContext) {
  console.log("attempting to read .blaim file");
  for (const ws of vscode.workspace.workspaceFolders || []) {
    const blaimContents = await vscode.workspace.fs.readFile(
      vscode.Uri.file(ws.uri.fsPath + "/.blaim"),
    );
    const blaimJsonStr = new TextDecoder().decode(blaimContents);
    const blaimJson: AcceptLogLine[] = JSON.parse(blaimJsonStr);
    addAcceptsFromBlaimFile(blaimJson);
  }
}

export function activate(context: vscode.ExtensionContext) {
  console.log("blaim-completion started");
  loadBlaimFile(context);
  activateDecorators(context);

  console.log(
    `writing accepted suggestions logs: export ACCEPT_LOG=${context.logUri.path}/accepted.suggestions.log`,
  );

  // This is the callback that VS Code invokes whenever the user *accepts* a suggestion.
  vscode.commands.registerCommand(acceptSuggestCommand, async (...args) => {
    const acceptEvent = args[0];
    const inferenceConfig = acceptEvent.inferenceConfig;

    const gitExtension =
      vscode.extensions.getExtension<GitExtension>("vscode.git")?.exports;
    const git = gitExtension?.getAPI(1);

    // Dig through the available repos to identify which one, if any,
    // owns the file the user is currently editing.
    for (const repo of git?.repositories || []) {
      if (acceptEvent.fileName.indexOf(repo?.rootUri.path) === 0) {
        const headCommit = repo?.state.HEAD;

        const acceptLogLine: AcceptLogLine = {
          fileName: acceptEvent.fileName.substring(
            repo?.rootUri.path.length + 1,
          ),
          position: acceptEvent.position,
          text: acceptEvent.text,
          headGitCommit: headCommit,
          inferenceConfig: inferenceConfig,
        };
        logAcceptedSuggestion(acceptLogLine);
      }
    }
  });

  const provider: vscode.InlineCompletionItemProvider = {
    async provideInlineCompletionItems(document, position, context, token) {
      if (token.isCancellationRequested) {
        console.log(`Canceled before AI completion.`);
        return;
      }
      return await lock.inLock(async () => {
        if (token.isCancellationRequested) {
          console.log(`Canceled before AI completion (inside lock).`);
          return;
        }

        const promptData = getPromptData(document, position);
        const prompt = formatPrompt(
          modelName,
          promptData.prefix,
          promptData.suffix,
        );
        const req: GenerateRequest & { stream: false } = {
          model: modelName,
          prompt: prompt.prompt,
          raw: true,
          stream: false,
          options: {
            stop: prompt.stop,
            num_predict: 20,
            temperature: 0.2,
          },
        };
        console.log("sending ollama.generate request...");
        const response = await ollama.generate(req);
        const text = response.response;
        console.log("ollama.generate response text length:", text.length);
        const ret: vscode.InlineCompletionList = {
          items: [
            {
              insertText: text,
              range: new vscode.Range(position, position),
              command: {
                title: "Accept generated code suggestion",
                command: acceptSuggestCommand,
                arguments: [
                  {
                    fileName: document.fileName,
                    fileVersionNumber: document.version,
                    position: position,
                    text: text,
                    inferenceConfig: {
                      modelName: req.model,
                      temperature: req.options?.temperature,
                      maxTokens: req.options?.num_predict,
                    },
                  },
                ],
              },
            },
          ],
        };
        return ret;
      });
    },

    handleDidShowCompletionItem(
      completionItem: vscode.InlineCompletionItem,
    ): void {
      //console.log("handleDidShowCompletionItem", completionItem);
    },

    /**
     * I still can't figure out what this callback is good for, but all the example code I've seen
     * implements it.  It gets called a lot, though it does not appear to get called when the
     * user actually accepts the inline suggestion.  For that, we have to implement the
     * acceptSuggestCommand thing above.
     */
    handleDidPartiallyAcceptCompletionItem(
      completionItem: vscode.InlineCompletionItem,
      info: any,
    ): void {
      // console.log(
      //   "handleDidPartiallyAcceptCompletionItem",
      //   completionItem,
      //   info,
      // );
    },
  };
  vscode.languages.registerInlineCompletionItemProvider(
    { pattern: "**", scheme: "file" },
    provider,
  );
}
