import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { Button, NavBar, Tabbar, TabbarItem, Swipe, SwipeItem, Search, Grid, GridItem, Card, Icon, Stepper, SubmitBar, Tag, List, PullRefresh, Cell, CellGroup, Field, Form, Popup, Picker, Toast, Dialog, ActionSheet, Uploader, Rate, Empty, Loading, Image as VanImage } from 'vant'
import 'vant/lib/index.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.use(Button)
app.use(NavBar)
app.use(Tabbar)
app.use(TabbarItem)
app.use(Swipe)
app.use(SwipeItem)
app.use(Search)
app.use(Grid)
app.use(GridItem)
app.use(Card)
app.use(Icon)
app.use(Stepper)
app.use(SubmitBar)
app.use(Tag)
app.use(List)
app.use(PullRefresh)
app.use(Cell)
app.use(CellGroup)
app.use(Field)
app.use(Form)
app.use(Popup)
app.use(Picker)
app.use(Toast)
app.use(Dialog)
app.use(ActionSheet)
app.use(Uploader)
app.use(Rate)
app.use(Empty)
app.use(Loading)
app.use(VanImage)

app.mount('#app')
