import * as vscode from 'vscode';
import { Range } from 'vscode';

console.log('got here');

export function activate(context: vscode.ExtensionContext) {
	console.log('inline-completions demo started');
	vscode.commands.registerCommand('demo-ext.command1', async (...args) => {
		console.log('command1: ' + JSON.stringify(args));
		//vscode.window.showInformationMessage('command1: ' + JSON.stringify(args));
	});

	const provider: vscode.InlineCompletionItemProvider = {
		async provideInlineCompletionItems(document, position, context, token) {
			console.log('provideInlineCompletionItems triggered');
			const regexp = /\/\/ \[(.+?),(.+?)\)(.*?):(.*)/;
			if (position.line <= 0) {
				return;
			}

			const result: vscode.InlineCompletionList = {
				items: [],
				commands: [],
			};

			let offset = 1;
			while (offset > 0) {
				if (position.line - offset < 0) {
					break;
				}
				
				const lineBefore = document.lineAt(position.line - offset).text;
				const matches = lineBefore.match(regexp);
				if (!matches) {
					break;
				}
				offset++;

				const start = matches[1];
				const startInt = parseInt(start, 10);
				const end = matches[2];
				const endInt =
					end === '*'
						? document.lineAt(position.line).text.length
						: parseInt(end, 10);
				const flags = matches[3];
				const completeBracketPairs = flags.includes('b');
				const isSnippet = flags.includes('s');
				const text = matches[4].replace(/\\n/g, '\n');

				result.items.push({
					insertText: isSnippet ? new vscode.SnippetString(text) : text,
					range: new Range(position.line, startInt, position.line, endInt),
					completeBracketPairs,
					command: {
						command: 'demo-ext.command1',
						title: 'My Inline Completion Demo Command',
						arguments: [{fileName: document.fileName, position, text}],
					}
				});
			}
			return result;
		},

		handleDidShowCompletionItem(completionItem: vscode.InlineCompletionItem): void {
			//console.log('handleDidShowCompletionItem', completionItem);
		},

		/**
		 * Is called when an inline completion item was accepted partially.
		 * @param acceptedLength The length of the substring of the inline completion that was accepted already.
		 */
		handleDidPartiallyAcceptCompletionItem(
			completionItem: vscode.InlineCompletionItem,
			info: any
		): void { 
			console.log('handleDidPartiallyAcceptCompletionItem', completionItem, info);
			const alreadyAcceptedPortion = completionItem.insertText.toString().substring(info);
			console.log('already accepted', alreadyAcceptedPortion);
		},
	};
	vscode.languages.registerInlineCompletionItemProvider({ pattern: '**', scheme: 'file' }, provider);
}
