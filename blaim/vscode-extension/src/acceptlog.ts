import * as vscode from 'vscode';

export interface AcceptLogLine {
  fileName: string,
  position: vscode.Position,
  text: string,
  headGitCommit: any,
  inferenceConfig: any,
}

export const accepts: AcceptLogLine[] = [];
