import React, { useState, useEffect } from "react";
import { type Agent } from "../../models/Agent";
import SidebarModal from "./SideBarModal";
import { Terminal, Info, Zap, History } from "lucide-react";
import CommandTimeline from "../timeline/CommandTimeLine.tsx";
import type { CommandEntry } from "../timeline/CommandTimeLine.tsx";
import CommandInput from "../common/CommandInput.tsx";

interface InteractAgentSideBarProps {
  isOpen: boolean;
  onClose: () => void;
  agent: Agent | null;
}

type TabOption = "terminal" | "timeline" | "extended" | "actions";

const InteractionAgentSideBar: React.FC<InteractAgentSideBarProps> = ({
  isOpen,
  onClose,
  agent,
}) => {
  const [activeTab, setActiveTab] = useState<TabOption>("terminal");
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);
  const [commands, setCommands] = useState<CommandEntry[]>([]);

  // Reset when agent changes
  useEffect(() => {
    if (agent) {
      setTerminalOutput([
        "[+] Connected to agent...",
        "[+] Awaiting instructions...",
      ]);
      setActiveTab("terminal");
    } else {
      setTerminalOutput([]);
    }
  }, [agent]);

  // Fetch existing commands for this agent
  useEffect(() => {
    if (!agent) return;
    const fetchCommands = async () => {
      try {
        const res = await fetch(`/api/v1/general/agents/${agent.id}/commands`);
        const data = await res.json();
        setCommands(
          data.map((cmd: any) => ({
            id: cmd.id,
            command: cmd.command,
            status: cmd.status,
            createdAt: cmd.created_at,
          }))
        );
      } catch (err) {
        console.error("Error fetching commands", err);
      }
    };
    fetchCommands();
  }, [agent]);

  // Poll for command results every 3s
  useEffect(() => {
    if (!agent) return;

    const interval = setInterval(async () => {
      for (const cmd of commands) {
        if (cmd.status !== "completed" && cmd.status !== "failed") {
          try {
            const res = await fetch(
              `/api/v1/general/commands/${cmd.id}/result`
            );
            const data = await res.json();
            setCommands((prev) =>
              prev.map((c) =>
                c.id === cmd.id
                  ? { ...c, status: data.status, result: data.output }
                  : c
              )
            );
          } catch (err) {
            console.error("Error fetching command result", err);
          }
        }
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [commands, agent]);

  const handleSendCommand = async (command: string, type: string = "shell") => {
    if (!agent) return;
    try {
      const res = await fetch(`/api/v1/general/queue/command/${agent.id}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ command, type }), // Matches Go RawCommand struct
      });

      if (!res.ok) {
        throw new Error(`Failed to send command: ${res.statusText}`);
      }

      const data = await res.json();

      setCommands((prev) => [
        {
          id: data.id,
          command,
          status: data.status,
          createdAt: new Date().toISOString(),
        },
        ...prev,
      ]);

      setTerminalOutput((prev) => [
        ...prev,
        `[>] Sending instruction: ${command} (${type})`,
        `[+] Awaiting response...`,
      ]);
    } catch (err) {
      console.error("Error sending command", err);
    }
  };

  return (
    <SidebarModal
      isOpen={isOpen}
      onClose={onClose}
      title={agent ? `Agent: ${agent.hostname}` : "No Agent Selected"}
    >
      {agent ? (
        <div className="space-y-6">
          {/* General Info */}
          <div className="space-y-2">
            <h3 className="text-sm uppercase text-gray-400 dark:text-gray-500 font-semibold">
              General Info
            </h3>
            <div className="text-gray-800 dark:text-gray-200 text-sm space-y-1">
              <p>
                <strong>Hostname:</strong> {agent.hostname}
              </p>
              <p>
                <strong>IP Address:</strong> {agent.ip}
              </p>
              <p>
                <strong>OS:</strong> {agent.os}
              </p>
              <p>
                <strong>Protocol:</strong> {agent.protocol}
              </p>
              <p>
                <strong>Status:</strong>{" "}
                {agent.isActive ? "Active" : "Inactive"}
              </p>
              <p>
                <strong>Last Seen:</strong>{" "}
                {new Date(agent.lastSeen).toLocaleString()}
              </p>
            </div>
          </div>

          {/* Tabs */}
          <div>
            <div className="flex border-b border-gray-200 dark:border-gray-700">
              {[
                {
                  key: "terminal",
                  label: "Terminal",
                  icon: <Terminal className="w-4 h-4 mr-2" />,
                },
                {
                  key: "timeline",
                  label: "Timeline",
                  icon: <History className="w-4 h-4 mr-2" />,
                },
                {
                  key: "extended",
                  label: "Extended Info",
                  icon: <Info className="w-4 h-4 mr-2" />,
                },
                {
                  key: "actions",
                  label: "Actions",
                  icon: <Zap className="w-4 h-4 mr-2" />,
                },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as TabOption)}
                  className={`flex items-center px-3 py-2 text-sm rounded-t-lg transition-colors
                    ${
                      activeTab === tab.key
                        ? "text-violet-600 dark:text-violet-400 bg-violet-100 dark:bg-violet-900/40 font-medium"
                        : "text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700/60"
                    }`}
                >
                  {tab.icon}
                  {tab.label}
                </button>
              ))}
            </div>

            {/* Tab Contents */}
            <div className="mt-4">
              {activeTab === "terminal" && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Terminal Output
                  </h3>
                  <div className="bg-black dark:bg-gray-900 text-green-400 font-mono text-sm p-3 rounded-lg h-64 overflow-y-auto border border-gray-700 dark:border-gray-600">
                    {terminalOutput.map((line, idx) => (
                      <div key={idx}>{line}</div>
                    ))}
                  </div>
                </div>
              )}

              {activeTab === "timeline" && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Command Timeline
                  </h3>
                  <CommandTimeline commands={commands} />

                  <CommandInput onSend={handleSendCommand} />
                </div>
              )}

              {activeTab === "extended" && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Extended Info
                  </h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Extended agent details will appear here (to be implemented).
                  </p>
                </div>
              )}

              {activeTab === "actions" && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2">
                    Quick Actions
                  </h3>
                  <button
                    type="button"
                    onClick={() => handleSendCommand("whoami", "shell")}
                    className="w-full bg-violet-600 text-white py-2.5 px-5 rounded-lg hover:bg-violet-700 focus:ring-4 focus:outline-none focus:ring-violet-300 font-medium text-sm text-center dark:bg-violet-500 dark:hover:bg-violet-600 dark:focus:ring-violet-800"
                  >
                    Run whoami
                  </button>

                  <button
                    type="button"
                    onClick={() => handleSendCommand("get-system-info", "intern")}
                    className="w-full mt-2 bg-indigo-600 text-white py-2.5 px-5 rounded-lg hover:bg-indigo-700 focus:ring-4 focus:outline-none focus:ring-indigo-300 font-medium text-sm text-center dark:bg-indigo-500 dark:hover:bg-indigo-600 dark:focus:ring-indigo-800"
                  >
                    Run internal system check
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      ) : (
        <p className="text-gray-500 dark:text-gray-400">No agent selected</p>
      )}
    </SidebarModal>
  );
};

export default InteractionAgentSideBar;
