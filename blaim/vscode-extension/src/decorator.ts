import * as vscode from "vscode";
import { getAcceptedSuggestionsForFile, AcceptLogLine } from "./acceptlog";

// this method is called when vs code is activated
export function activateDecorators(context: vscode.ExtensionContext) {
  console.log("decorator sample is activated");

  const toggleCallback = async function name(params: any) {
    console.log("toggle", params);
  };

  const toggle = vscode.commands.registerCommand(
    "coverage-gutters.toggleCoverage",
    toggleCallback,
  );

  context.subscriptions.push(toggle);

  let timeout: string | number | NodeJS.Timeout | undefined = undefined;

  // const gutterBlameAnnotationType = vscode.window.createTextEditorDecorationType({
  // 	rangeBehavior: vscode.DecorationRangeBehavior.ClosedClosed,
  //   isWholeLine: true,
  //   gutterIconSize: 'contain',
  //   before: {
  //     color: 'gray',
  //     width: '100px',
  //     height: '100%',
  //     margin: '0 26px -1px 0',
  //     contentText: ''
  //   },
  // });

  const inlineDecorationType = vscode.window.createTextEditorDecorationType({
    rangeBehavior: vscode.DecorationRangeBehavior.ClosedClosed,
    textDecoration: "opacity:0.5",
    overviewRulerColor: "blue",
    overviewRulerLane: vscode.OverviewRulerLane.Right,
    light: {
      // this color will be used in light color themes
      backgroundColor: "lightgray",
      //borderColor: 'darkblue'
    },
    dark: {
      //borderRadius: '8px',
      // this color will be used in dark color themes
      //borderColor: 'lightblue'
      backgroundColor: "darkgreen",
    },
  });

  let activeEditor = vscode.window.activeTextEditor;

  function findFirstDiffPos(a: string, b: string) {
    let i = 0;
    if (a === b) return -1;
    while (a[i] === b[i]) i++;
    //console.log('findFirstDiffPos', a, b, i);
    return i;
  }

  function updateDecorations() {
    if (!activeEditor) {
      return;
    }
    const text = activeEditor.document.getText();
    const inlineDecorations: vscode.DecorationOptions[] = [];
    const gutterBlameAnnotations: vscode.DecorationOptions[] = [];

    // Inline annotations to hilight ranges of AI-generated code.
    // This is currently very broken and does not work well at all.
    const accepts = getAcceptedSuggestionsForFile(
      activeEditor.document.fileName,
    );
    for (let i = 0; i < accepts.length; i++) {
      const acceptLine = accepts[i];
      const acceptedLines = acceptLine.text.split("\n");
      const endPos = new vscode.Position(
        acceptLine.position.line + acceptedLines.length,
        acceptedLines?.pop()?.length || 0,
      );
      const range = new vscode.Range(acceptLine.position, endPos);

      // At this point, we have the originally logged suggestion text
      // and the position where it was inserted in the original file.
      //
      // However, this is not enough information to let us higlight the
      // *current* location of the suggested text in the file, if it is still
      // there at all. For instance:
      // - The user may accept a suggestion, and subsequently edit or remove
      //   the inserted text altogether.
      // - The user may accept a suggestion, and then make edits inserting new
      //   lines at a previous line in the file, such that the logged
      // . line number of the inserted text is no longer accurate.
      //
      // This code does nothing fancy to find fuzzy matches etc.  It looks
      // for an exact match at the original location, and as a fallback it
      // looks for exact matches at other locations in the document.

      // Look for exact matches of the accepted text at the original location:
      const actualTextAfterEdits = activeEditor.document.getText(range);
      if (actualTextAfterEdits == acceptLine.text) {
        const decoration: vscode.DecorationOptions = {
          range: range,
          hoverMessage: new vscode.MarkdownString(
            "## This is AI-generated code:\n```" +
              JSON.stringify(acceptLine.inferenceConfig) +
              "```\n" +
              `originally at ${JSON.stringify(acceptLine.position)}\n` +
              "## Raw text from the model:\n```\n" +
              acceptLine.text +
              "\n```\n",
          ),
        };
        inlineDecorations.push(decoration);
      } else {
        // If we didn't the exact text at the original location, try this as a back-up:
        // Look for exact matches of the accepted text at other locations in the file:
        let acceptTextIndexInDoc = activeEditor.document
          .getText()
          .indexOf(acceptLine.text);
        const ranges: vscode.Range[] = [];
        while (acceptTextIndexInDoc != -1) {
          ranges.push(
            new vscode.Range(
              activeEditor.document.positionAt(acceptTextIndexInDoc),
              activeEditor.document.positionAt(
                acceptTextIndexInDoc + acceptLine.text.length,
              ),
            ),
          );
          acceptTextIndexInDoc = activeEditor.document
            .getText()
            .indexOf(acceptLine.text, acceptTextIndexInDoc + 1);
        }
        for (let i = 0; i < ranges.length; i++) {
          const decoration: vscode.DecorationOptions = {
            range: ranges[i],
            hoverMessage: new vscode.MarkdownString(
              "## This is AI-generated code:\n```" +
                JSON.stringify(acceptLine.inferenceConfig) +
                "```\n" +
                `originally at ${JSON.stringify(acceptLine.position)}\n` +
                "## Raw text from the model:\n```\n" +
                acceptLine.text +
                "\n```\n",
            ),
          };
          inlineDecorations.push(decoration);
        }
      }

      // Line-by-line gutter annotations
      const textLines = text.split("\n");
      for (let lineNumber = 0; lineNumber < textLines.length; lineNumber++) {
        const linePos = new vscode.Position(lineNumber, 0);
        const gutterAnnotation: vscode.DecorationOptions = {
          range: new vscode.Range(linePos, linePos),
          //hoverMessage: 'blaim hover message',
          renderOptions: {
            before: {
              contentText: `blaim line ${lineNumber}`,
            },
          },
        };
        gutterBlameAnnotations.push(gutterAnnotation);
      }

      activeEditor.setDecorations(inlineDecorationType, []);
      activeEditor.setDecorations(inlineDecorationType, inlineDecorations);
      //    activeEditor.setDecorations(gutterBlameAnnotationType, gutterBlameAnnotations);
    }
  }

  function triggerUpdateDecorations(throttle = false) {
    if (timeout) {
      clearTimeout(timeout);
      timeout = undefined;
    }
    if (throttle) {
      timeout = setTimeout(updateDecorations, 500);
    } else {
      updateDecorations();
    }
  }

  if (activeEditor) {
    triggerUpdateDecorations();
  }

  vscode.window.onDidChangeActiveTextEditor(
    (editor) => {
      activeEditor = editor;
      if (editor) {
        triggerUpdateDecorations();
      }
    },
    null,
    context.subscriptions,
  );

  vscode.workspace.onDidChangeTextDocument(
    (event) => {
      if (activeEditor && event.document === activeEditor.document) {
        triggerUpdateDecorations(true);
      }
    },
    null,
    context.subscriptions,
  );
}
