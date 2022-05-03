import {createRouter, createWebHistory} from "vue-router"
import States from "./components/States";
import ServerInfo from "./components/ServerInfo";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: "/",
            component: States
        },
        {
            path: "/servers/:id",
            component: ServerInfo,
            props: true
        }
    ]
})

export default router