import * as vscode from "vscode";

export interface AcceptLogLine {
  fileName: string;
  position: vscode.Position;
  text: string;
  headGitCommit: any;
  inferenceConfig: any;
}

const accepts: AcceptLogLine[] = [];

const acceptLogger = vscode.window.createOutputChannel("accepted.suggestions", {
  log: true,
});

export function logAcceptedSuggestion(logLine: AcceptLogLine) {
  acceptLogger.appendLine(JSON.stringify(logLine));

  accepts.push(logLine);
}

export function getAcceptedSuggestionsForFile(path: string): AcceptLogLine[] {
  const ret: AcceptLogLine[] = [];
  for (let i = 0; i < accepts.length; i++) {
    const acceptLine = accepts[i];
    if (path.indexOf(acceptLine.fileName) != -1) {
      ret.push(acceptLine);
    }
  }

  return ret;
}
