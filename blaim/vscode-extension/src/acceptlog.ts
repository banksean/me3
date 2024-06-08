import { appendFile } from "fs";
import * as vscode from "vscode";

// TODO: reconsile this with https://github.com/microsoft/vscode-extension-samples/blob/main/telemetry-sample/src/extension.ts
// - there's probably a more proper way to do this.

export interface AcceptLogLine {
  fileName: string;
  position: vscode.Position;
  text: string;
  headGitCommit: any;
  inferenceConfig: any;
}

const accepts: AcceptLogLine[] = [];
let acceptsFromBlaimFile: AcceptLogLine[] = [];
export function addAcceptsFromBlaimFile(accepts: AcceptLogLine[]) {
  acceptsFromBlaimFile = accepts;
}

const acceptLogger = vscode.window.createOutputChannel("accepted.suggestions", {
  log: true,
});

export function logAcceptedSuggestion(logLine: AcceptLogLine) {
  acceptLogger.appendLine(JSON.stringify(logLine));
  accepts.push(logLine);
  console.log("accepts so far: ", accepts.length);
}

export function getAcceptedSuggestionsForFile(path: string): AcceptLogLine[] {
  const ret: AcceptLogLine[] = [];

  // Add any annotations for this file from the latest .blaim contents.
  for (let i = 0; i < acceptsFromBlaimFile.length; i++) {
    const acceptLine = acceptsFromBlaimFile[i];
    if (path.indexOf(acceptLine.fileName) != -1) {
      ret.push(acceptLine);
    }
  }

  // Add any in-memory annotaitons we have in the current session.
  for (let i = 0; i < accepts.length; i++) {
    const acceptLine = accepts[i];
    if (path.indexOf(acceptLine.fileName) != -1) {
      ret.push(acceptLine);
    }
  }

  return ret;
}
