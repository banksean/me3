import * as vscode from 'vscode';

// this method is called when vs code is activated
export function activateDecorators(context: vscode.ExtensionContext) {

	console.log('decorator sample is activated');

  const toggleCallback = async function name(params: any) {
    console.log('toggle', params);
  };
  
  const toggle = vscode.commands.registerCommand(
    "coverage-gutters.toggleCoverage",
    toggleCallback,
//    gutters.toggleCoverageForActiveFile.bind(gutters),
  );

  context.subscriptions.push(toggle);

  let timeout: string | number | NodeJS.Timeout | undefined = undefined;

	const gutterBlameAnnotationType = vscode.window.createTextEditorDecorationType({
		rangeBehavior: vscode.DecorationRangeBehavior.OpenOpen,
    isWholeLine: true,
    gutterIconSize: 'contain',
    before: {
      color: 'lightgray',
//      backgroundColor: 'blue',
      width: '100px',
      height: '100%',
      margin: '0 26px -1px 0',
      contentText: ''
    },
	});

  // const gutterBlameHighlightType = vscode.window.createTextEditorDecorationType({
  // //  gutterIconPath: gutterHighlightUri,
  //   gutterIconSize: 'contain',
  //   isWholeLine: true,
  //   overviewRulerLane: vscode.OverviewRulerLane.Full,
  //   backgroundColor: 'green',
  //   overviewRulerColor: 'magenta'
  // });

  const smallNumberDecorationType = vscode.window.createTextEditorDecorationType({
    borderWidth: '1px',
    borderStyle: 'solid',
    overviewRulerColor: 'blue',
    overviewRulerLane: vscode.OverviewRulerLane.Right,
    light: {
        // this color will be used in light color themes
        borderColor: 'darkblue'
    },
    dark: {
        // this color will be used in dark color themes
        borderColor: 'lightblue'
    }
  });


	// create a decorator type that we use to decorate large numbers
	const largeNumberDecorationType = vscode.window.createTextEditorDecorationType({
		cursor: 'crosshair',
		// use a themable color. See package.json for the declaration and default values.
		backgroundColor: { id: 'myextension.largeNumberBackground' }
	});

	let activeEditor = vscode.window.activeTextEditor;

	function updateDecorations() {
		if (!activeEditor) {
			return;
		}
		const regEx = /\d+/g;
		const text = activeEditor.document.getText();
		const smallNumbers: vscode.DecorationOptions[] = [];
		const largeNumbers: vscode.DecorationOptions[] = [];
    const gutterBlameAnnotations: vscode.DecorationOptions[] = [];
    //const gutterHighlightAnnotations: vscode.DecorationOptions[] = [];

		let match;
		while ((match = regEx.exec(text))) {
			const startPos = activeEditor.document.positionAt(match.index);
			const endPos = activeEditor.document.positionAt(match.index + match[0].length);
			const decoration = { range: new vscode.Range(startPos, endPos), hoverMessage: 'Number **' + match[0] + '**' };
			if (match[0].length < 3) {
				smallNumbers.push(decoration);
			} else {
				largeNumbers.push(decoration);
			}
    }
    const textLines = text.split('\n');
    for (let lineNumber = 0; lineNumber < textLines.length; lineNumber++) {
      const linePos = new vscode.Position(lineNumber, 0);
      const gutterAnnotation: vscode.DecorationOptions = {
        range: new vscode.Range(linePos, linePos), hoverMessage: 'blaim hover message',
        renderOptions: {
          before: {
            contentText: `blaim line ${lineNumber}`,
          }
        }
        // renderOptions: {
        //   before: {
        //     color: 'red',
        //     backgroundColor: 'blue',
        //     width: '100px',
        //     height: '100%',
        //     margin: '0 26px -1px 0',
        //     contentText: 'gutter text'
        //   },
        //},
      };
      gutterBlameAnnotations.push(gutterAnnotation);
		}
    
		activeEditor.setDecorations(smallNumberDecorationType, smallNumbers);
		activeEditor.setDecorations(largeNumberDecorationType, largeNumbers);

    activeEditor.setDecorations(gutterBlameAnnotationType, gutterBlameAnnotations);
    //activeEditor.setDecorations(gutterBlameHighlightType, gutterHighlightAnnotations);
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