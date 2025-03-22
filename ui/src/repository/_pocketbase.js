import PocketBase from 'pocketbase';

let pb = null;
export const getPocketBase = () => {
  if (pb) return pb;
  pb = new PocketBase(import.meta.env.VITE_APP_BASE_API);
  return pb;
};

export const COLLECTION_NAME_ADMIN = "_superusers";
export const COLLECTION_NAME_ACCESS = "access";
export const COLLECTION_NAME_SSH = "ssh";
export const COLLECTION_NAME_SETTINGS = "settings";
export const COLLECTION_NAME_PORT_FORWARDING = "port_forwarding";
////////////
export const COLLECTION_NAME_CONN_HISTORY = "conn_history";
export const COLLECTION_NAME_FOLDER = "folder";
export const COLLECTION_NAME_SNIPPET = "snippet";
export const COLLECTION_NAME_CMD_HISTORY = "cmd_history";
