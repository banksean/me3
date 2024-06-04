import * as vscode from 'vscode';
import {accepts, AcceptLogLine } from './acceptlog';

// this method is called when vs code is activated
export function activateDecorators(context: vscode.ExtensionContext) {

	console.log('decorator sample is activated');

  const toggleCallback = async function name(params: any) {
    console.log('toggle', params);
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
    overviewRulerColor: 'blue',
    overviewRulerLane: vscode.OverviewRulerLane.Right,
    light: {
        // this color will be used in light color themes
        backgroundColor: 'lightgray',
        //borderColor: 'darkblue'
    },
    dark: {
        //borderRadius: '8px',
        // this color will be used in dark color themes
        //borderColor: 'lightblue'
        backgroundColor: 'midnightblue'
    }
  });

	let activeEditor = vscode.window.activeTextEditor;

	function updateDecorations() {
		if (!activeEditor) {
			return;
		}
		const text = activeEditor.document.getText();
		const inlineDecorations: vscode.DecorationOptions[] = [];
    const gutterBlameAnnotations: vscode.DecorationOptions[] = [];

    for (let i =0; i<accepts.length; i++) {
      const acceptLine = accepts[i];
      if (activeEditor.document.fileName.indexOf(acceptLine.fileName) != -1) {
        const acceptedLines = acceptLine.text.split('\n');
        const endPos = new vscode.Position(acceptLine.position.line + acceptedLines.length, acceptedLines?.pop()?.length || 0);
        const decoration:vscode.DecorationOptions = { 
          range: new vscode.Range(acceptLine.position, endPos), 
          hoverMessage: new vscode.MarkdownString("## This is AI-generated code:\n```" + JSON.stringify(acceptLine.inferenceConfig) + "```\n"  + `at ${JSON.stringify(acceptLine.position)}\n`+ "## Raw text from the model:\n```\n" + acceptLine.text + "\n```\n"),
        };
        inlineDecorations.push(decoration);
      }
    }
    // Line-by-line gutter annotations
    const textLines = text.split('\n');
    for (let lineNumber = 0; lineNumber < textLines.length; lineNumber++) {
      const linePos = new vscode.Position(lineNumber, 0);
      const gutterAnnotation: vscode.DecorationOptions = {
        range: new vscode.Range(linePos, linePos), 
        //hoverMessage: 'blaim hover message',
        renderOptions: {
          before: {
            contentText: `blaim line ${lineNumber}`,
          }
        }
      };
      gutterBlameAnnotations.push(gutterAnnotation);
		}
    
		activeEditor.setDecorations(inlineDecorationType, inlineDecorations);
    //    activeEditor.setDecorations(gutterBlameAnnotationType, gutterBlameAnnotations);
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

	vscode.window.onDidChangeActiveTextEditor(editor => {
		activeEditor = editor;
		if (editor) {
			triggerUpdateDecorations();
		}
	}, null, context.subscriptions);

	vscode.workspace.onDidChangeTextDocument(event => {
		if (activeEditor && event.document === activeEditor.document) {
			triggerUpdateDecorations(true);
		}
	}, null, context.subscriptions);

  
}