import Vue from "vue";
import {
  Cell,
  Button,
  List,
  Loading,
  CellGroup,
  Row,
  Col,
  Panel,
  Field,
  Picker,
  Popup,
  Toast,
  NumberKeyboard,
  Dialog,
  Tab,
  Tabs,
  Image,
  Icon,
  SwipeCell,
  Switch
} from "vant";
import { Locale } from "vant";
import enUS from "vant/lib/locale/lang/en-US";
import zhCN from "vant/lib/locale/lang/zh-CN";
import i18n from "@/i18n";
import "vant/lib/index.css";

Locale.use("en-US", enUS);
if (i18n.locale === "en") {
  Locale.use("en-US", enUS);
} else {
  Locale.use("zh-CN", zhCN);
}

Vue.use(Toast);
Vue.use(Cell);
Vue.use(List);
Vue.use(Button);
Vue.use(Loading);
Vue.use(CellGroup).use(Cell);
Vue.use(Row).use(Col);
Vue.use(Panel);
Vue.use(Field);
Vue.use(Picker);
Vue.use(Popup);
Vue.use(NumberKeyboard);
Vue.use(Dialog);
Vue.use(Tab);
Vue.use(Tabs);
Vue.use(Image);
Vue.use(Icon);
Vue.use(Switch);
Vue.use(SwipeCell);
