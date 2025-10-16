import React from "react";
import { CheckCircle2, Loader2, XCircle } from "lucide-react";

const SUPPRESSED_COMMANDS = ["return_results"];

export interface CommandEntry {
  id: string;
  command: string;
  status: string;
  createdAt: string;
  result?: string;
}

interface CommandTimelineProps {
  commands: CommandEntry[];
}

const statusMap: Record<string, string> = {
  cs_pndg: "Pending",
  cs_rng: "Running",
  cs_cmpltd: "Completed",
  cs_fld: "Failed",
  cs_clld: "Cancelled",
  cs_tmt: "Timeout",
  cs_ack: "Acknowledged",
  c_sent: "Sent to Client",
  s_sent: "Sent to Server",
  c_received: "Received from Client",
  s_received: "Received from Server",
};

const CommandTimeline: React.FC<CommandTimelineProps> = ({ commands }) => {
  const visibleCommands = commands.filter(
    (cmd) => !SUPPRESSED_COMMANDS.includes(cmd.command)
  );

  if (visibleCommands.length === 0) {
    return (
      <p className="text-sm text-gray-500 dark:text-gray-400">
        No commands executed yet.
      </p>
    );
  }

  const getStatusColor = (status: string): string => {
    if (status === "cs_cmpltd") return "text-green-500";
    if (status === "cs_fld") return "text-red-500";
    return "text-orange-500";
  };

  const getBorderColor = (status: string): string => {
    if (status === "cs_cmpltd") return "border-green-500";
    if (status === "cs_fld") return "border-red-500";
    return "border-orange-500";
  };

  return (
    <div className="space-y-4 max-h-72 overflow-y-auto pr-2">
      {visibleCommands.map((cmd) => (
        <div
          key={cmd.id}
          className={`border-l-2 pl-4 relative ${getBorderColor(cmd.status)}`}
        >
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-gray-800 dark:text-gray-200">
              {cmd.command}
            </p>
            <span className="text-xs text-gray-500 dark:text-gray-400">
              {new Date(cmd.createdAt).toLocaleTimeString()}
            </span>
          </div>
          <div className="flex items-center space-x-2 mt-1">
            {cmd.status === "cs_cmpltd" && (
              <CheckCircle2 className="w-4 h-4 text-green-500" />
            )}
            {cmd.status === "cs_fld" && (
              <XCircle className="w-4 h-4 text-red-500" />
            )}
            {cmd.status !== "cs_cmpltd" && cmd.status !== "cs_fld" && (
              <Loader2 className="w-4 h-4 text-orange-500 animate-spin" />
            )}
            <span
              className={`text-xs font-medium ${getStatusColor(cmd.status)}`}
            >
              {statusMap[cmd.status] || cmd.status}
            </span>
          </div>
          {cmd.result && (
            <pre className="bg-gray-900 text-green-400 font-mono text-xs mt-2 p-2 rounded-lg whitespace-pre-wrap border border-gray-700 dark:border-gray-600">
              {cmd.result}
            </pre>
          )}
        </div>
      ))}
    </div>
  );
};

export default CommandTimeline;