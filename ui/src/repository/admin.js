import { COLLECTION_NAME_ADMIN, getPocketBase } from "./_pocketbase";

export const authWithPassword = (username, password) => {
  return getPocketBase().collection(COLLECTION_NAME_ADMIN).authWithPassword(username, password);
};

export const getAuthStore = () => {
  return getPocketBase().authStore;
};

export const save = (data) => {
  return getPocketBase()
    .collection(COLLECTION_NAME_ADMIN)
    .update(getAuthStore().record?.id || "", data);
};
