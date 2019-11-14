import Vue from "vue";
import VueRouter from "vue-router";
import Home from "./pages/Home";
import TestAuth from "./pages/TestAuth";
import Pay from "./pages/Pay";
import PayWxQr from "./pages/PayWxQr";
import PreparePacket from "./pages/PreparePacket";
import Packet from "./pages/Packet";
import SendReward from "./pages/rewards/Send";
import RewardRank from "./pages/rewards/Rank";
import Members from "./pages/Members";
import Coupons from "./pages/Coupons";
import Invitations from "./pages/Invitations";
import Messages from "./pages/Messages";
import WxPay from "./pages/WxPay";
import WxPayDone from "./pages/WxPayDone";
import GrantToken from "./pages/GrantToken";
import { ROUTER_MODE } from "@/constants.js";
import InvitationDetails from "./pages/InvitationDetails";
import InvitationEntry from "./pages/InvitationEntry";
import EditAnnouncement from "./pages/settings/EditAnnouncement";
import GroupMode from "./pages/settings/GroupMode";

Vue.use(VueRouter);

const routes = [
  { path: "/", component: Home },
  { path: "/pay", component: Pay },
  { path: "/pay/wxqr", component: PayWxQr },
  { path: "/packets/prepare", component: PreparePacket },
  { path: "/packets/:id", component: Packet },
  { path: "/coupons/", component: Coupons },
  { path: "/rewards/send", component: SendReward },
  { path: "/rewards/rank", component: RewardRank },
  { path: "/auth", component: TestAuth },
  { path: "/invitations/", component: Invitations },
  { path: "/grant-token/", component: GrantToken },
  // special route for wechat, no auth required.
  { path: "/wxpay", component: WxPay },
  { path: "/wxpay/done", component: WxPayDone },
  { path: "/invitation/details", component: InvitationDetails },
  { path: "/invitation/entry", component: InvitationEntry },
  { path: "/settings/announcement/edit", component: EditAnnouncement },
  { path: "/settings/mode", component: GroupMode },
  { path: "/members/", component: Members },
  { path: "/messages/", component: Messages }
];

const router = new VueRouter({
  mode: ROUTER_MODE,
  routes // short for `routes: routes`
});

export default router;
