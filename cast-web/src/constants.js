export const isProduction = process.env.NODE_ENV === "production";

export const CLIENT_ID = process.env.VUE_APP_CLIENT_ID;

export const WEB_ROOT = process.env.VUE_APP_WEB_ROOT;

export const BASE_URL = process.env.VUE_APP_API_ROOT;

export const WS_BASE_URL = process.env.VUE_APP_WS_ROOT;

export const ROUTER_MODE = process.env.VUE_APP_ROUTER_MODE || "hash";

export const COLORS = {
  NAV_COLOR: process.env.VUE_APP_NAV_COLOR
}

export const MIXIN_HOST = "https://api.mixin.one";

export const EOS_ASSET_ID = "6cfe566e-4aad-470b-8c9a-2fd35b49c68d";

export const SUP_MESSAGE_CAT = ['PLAIN_TEXT', 'PLAIN_IMAGE', 'PLAIN_VIDEO', 'PLAIN_AUDIO', 'PLAIN_LIVE', 'PLAIN_DATA']

export const SOCKET_STATE = {
  DISCONNECT: 'disconnect',
  CONNECTED: 'connected',
  CONNECTING: 'connecting'
}
