import React from "react";
import { CheckCircle2, Loader2, XCircle } from "lucide-react";

const SUPPRESSED_COMMANDS = ["return_results"];

export interface CommandEntry {
  id: string;
  command: string;
  status: string;
  createdAt: string;
  result?: string;
  result_type?: string;
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

  const formatResult = (result: string, resultType?: string): React.ReactNode => {
    if (!result) return null;

    try {
      const parsed = JSON.parse(result);
      
      switch (resultType) {
        case "list":
          // Format as bullet list
          if (Array.isArray(parsed)) {
            return (
              <ul className="list-disc list-inside space-y-1">
                {parsed.map((item, idx) => (
                  <li key={idx} className="text-green-400">{item}</li>
                ))}
              </ul>
            );
          }
          break;
        
        case "key_value":
        case "structured":
          // Format as key-value pairs
          if (typeof parsed === "object" && !Array.isArray(parsed)) {
            return (
              <div className="space-y-1">
                {Object.entries(parsed).map(([key, value]) => (
                  <div key={key} className="flex gap-2">
                    <span className="text-blue-400 font-semibold">{key}:</span>
                    <span className="text-green-400">
                      {typeof value === "object" ? JSON.stringify(value, null, 2) : String(value)}
                    </span>
                  </div>
                ))}
              </div>
            );
          }
          break;
        
        case "json":
          // Pretty print JSON
          return (
            <pre className="text-green-400 whitespace-pre-wrap">
              {JSON.stringify(parsed, null, 2)}
            </pre>
          );
        
        case "error":
          // Highlight errors in red
          return (
            <div className="text-red-400 font-medium">
              {typeof parsed === "string" ? parsed : JSON.stringify(parsed, null, 2)}
            </div>
          );
        
        case "table":
        case "text":
        default:
          // Plain text or table - display as-is
          return (
            <span className="text-green-400 whitespace-pre-wrap">
              {typeof parsed === "string" ? parsed : JSON.stringify(parsed, null, 2)}
            </span>
          );
      }
    } catch {
      // If not valid JSON, display as plain text
      return <span className="text-green-400 whitespace-pre-wrap">{result}</span>;
    }
    
    return <span className="text-green-400 whitespace-pre-wrap">{result}</span>;
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
            <div className="bg-gray-900 font-mono text-xs mt-2 p-2 rounded-lg border border-gray-700 dark:border-gray-600">
              {formatResult(cmd.result, cmd.result_type)}
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default CommandTimeline;