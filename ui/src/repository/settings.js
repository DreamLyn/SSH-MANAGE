
import { COLLECTION_NAME_SETTINGS, getPocketBase } from "./_pocketbase";

export const get = async (name) => {
  try {
    const resp = await getPocketBase().collection(COLLECTION_NAME_SETTINGS).getFirstListItem(`name='${name}'`, {
      requestKey: null,
    });
    return resp;
  } catch (err) {
    if (err instanceof ClientResponseError && err.status === 404) {
      return {
        name: name,
        content: {},
      };
    }
    throw err;
  }
};

export const save = async (record) => {
  if (record.id) {
    return await getPocketBase().collection(COLLECTION_NAME_SETTINGS).update(record.id, record);
  }
  return await getPocketBase().collection(COLLECTION_NAME_SETTINGS).create(record);
};


export const subscribe = async (id, cb) => {
  return getPocketBase().collection(COLLECTION_NAME_SETTINGS).subscribe(id, cb);
};

export const unsubscribe = async (id) => {
  return getPocketBase().collection(COLLECTION_NAME_SETTINGS).unsubscribe(id);
};
