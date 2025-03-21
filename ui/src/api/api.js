import { ClientResponseError } from "pocketbase";

import { getPocketBase } from "@/repository/_pocketbase";


export const connectForwarding = async (forwardingId) => {
  const pb = getPocketBase();

  const resp = await pb.send(`/forwarding/${encodeURIComponent(forwardingId)}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      
    },
  });

  if (resp.code != 0) {
    console.log(resp)
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};
