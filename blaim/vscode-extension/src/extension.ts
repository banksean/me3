import * as vscode from "vscode";
import ollama from "ollama";
import { GenerateRequest } from "ollama";
import type { GitExtension } from "./git";

const modelName = "codellama";
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

function formatPrompt(prefix: string, suffix: string) {
  return {
    prompt: `<PRE> ${prefix} <SUF> ${suffix} <MID>`,
    stop: [`<END>`, `<EOD>`, `<EOT>`],
  };
}

export function activate(context: vscode.ExtensionContext) {
  console.log("blaim-completion started");
  const acceptLogger = vscode.window.createOutputChannel(
    "accepted.suggestions",
    { log: true },
  );
  console.log(
    `writing accepted suggestions logs: export ACCEPT_LOG=${context.logUri.path}/accepted.suggestions.log`,
  );

  vscode.commands.registerCommand(acceptSuggestCommand, async (...args) => {
    const acceptEvent = args[0];
    const inferenceCofig = acceptEvent.inferenceConfig;

    const gitExtension =
      vscode.extensions.getExtension<GitExtension>("vscode.git")?.exports;
    const git = gitExtension?.getAPI(1);

    for (const repo of git?.repositories || []) {
      if (acceptEvent.fileName.indexOf(repo?.rootUri.path) === 0) {
        const headCommit = repo?.state.HEAD;
        acceptLogger.appendLine(
          JSON.stringify({
            fileName: acceptEvent.fileName.substring(
              repo?.rootUri.path.length + 1,
            ),
            position: acceptEvent.position,
            text: acceptEvent.text,
            headGitCommit: headCommit,
            inferenceConfig: inferenceCofig,
          }),
        );
      }
    }
  });

  const provider: vscode.InlineCompletionItemProvider = {
    async provideInlineCompletionItems(document, position, context, token) {
      const promptData = getPromptData(document, position);
      const prompt = formatPrompt(promptData.prefix, promptData.suffix);
      const req: GenerateRequest & { stream: false } = {
        model: modelName,
        prompt: prompt.prompt,
        raw: true,
        stream: false,
        options: {
          stop: prompt.stop,
          num_predict: 5,
          temperature: 0.2,
        },
      };
      console.log("sending ollama.generate request...");
      const response = await ollama.generate(req);
      const text = response.response;
      console.log("ollama.generate response text length:", text.length);
			console.log('position:', position);
      const ret: vscode.InlineCompletionList = {
        items: [
          {
            insertText: text,
            range: new vscode.Range(position, position),
            command: {
              title: "Accept Llama suggestion",
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
    },

    handleDidShowCompletionItem(
      completionItem: vscode.InlineCompletionItem,
    ): void {
      console.log("handleDidShowCompletionItem", completionItem);
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
      console.log(
        "handleDidPartiallyAcceptCompletionItem",
        completionItem,
        info,
      );
    },
  };
  vscode.languages.registerInlineCompletionItemProvider(
    { pattern: "**", scheme: "file" },
    provider,
  );
}
