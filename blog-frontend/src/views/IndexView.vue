<script setup>
import {get, logout} from '@/net'
import router from "@/router";
import {useStore} from "@/store";
import {reactive, ref} from "vue";
import {
    Back,
    Bell,
    ChatDotSquare, Check, Collection, DataBoard, DataLine,
    Document, Files, Football,
    Location, Lock, Message, Monitor,
    Notification, Operation,
    Position,
    School, Search, Sunny,
    Umbrella, User, View
} from "@element-plus/icons-vue";

const store = useStore()
const loading = ref(true)

const searchInput = reactive({
    type: '1',
    text: ''
})

const notification = ref([])

get('/api/user/info', (data) => {
    store.user = data
    loading.value = false
})

function userLogout() {
    logout(() => router.push("/"))
}

</script>

<template>
    <div class="main-content" v-loading="loading" element-loading-text="正在进入，请稍后...">
        <el-container style="height: 100%" v-if="!loading">
            <el-header class="main-content-header">
                <el-image class="logo" src="https://element-plus.org/images/element-plus-logo.svg"></el-image>
<!--                <el-image class="logo" src="https://s2.loli.net/2024/05/09/Hhv6IFetQUcRmG7.png"></el-image>-->
                <div style="flex: 1" class="user-info">
                    <div class="profile">
                        <div>{{ store.user.username }}</div>
                        <div>{{ store.user.email }}</div>
                    </div>
                    <el-dropdown>
                        <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"/>
<!--                        <el-avatar :src="store.avatarUrl"/>-->
                        <template #dropdown>
                            <el-dropdown-item>
                                <el-icon>
                                    <Operation/>
                                </el-icon>
                                个人设置
                            </el-dropdown-item>
                            <el-dropdown-item @click="userLogout" divided>
                                <el-icon>
                                    <Back/>
                                </el-icon>
                                退出登录
                            </el-dropdown-item>
                        </template>
                    </el-dropdown>
                </div>
            </el-header>
            <el-container>
                <el-aside width="230px">
                    <el-scrollbar style="height: calc(100vh - 55px)">
                        <el-menu
                            router
                            :default-active="$route.path"
                            :default-openeds="['1', '2', '3']"
                            style="min-height: calc(100vh - 55px)">
                            <el-sub-menu index="1">
                                <template #title>
                                    <el-icon>
                                        <Location/>
                                    </el-icon>
                                    <span><b>校园论坛</b></span>
                                </template>
                                <el-menu-item index="/index">
                                    <template #title>
                                        <el-icon>
                                            <ChatDotSquare/>
                                        </el-icon>
                                        帖子广场
                                    </template>
                                </el-menu-item>
                            </el-sub-menu>
                            <el-sub-menu index=2>
                                <template #title>
                                    <el-icon>
                                        <Operation/>
                                    </el-icon>
                                    <span><b>个人设置</b></span>
                                </template>
                                <el-menu-item index="/index/user-setting">
                                    <template #title>
                                        <el-icon>
                                            <User/>
                                        </el-icon>
                                        个人信息设置
                                    </template>
                                </el-menu-item>
                            </el-sub-menu>
                        </el-menu>
                    </el-scrollbar>
                </el-aside>
                <el-main class="main-content-page">
                    <el-scrollbar style="height: calc(100vh - 55px)">
                        <router-view v-slot="{ Component }">
                            <transition name="el-fade-in-linear" mode="out-in">
                                <component :is="Component" style="height: 100%"/>
                            </transition>
                        </router-view>
                    </el-scrollbar>
                </el-main>
            </el-container>
        </el-container>
    </div>
</template>

<style lang="less" scoped>
.notification-item {
    transition: .3s;
    &:hover {
        cursor: pointer;
        opacity: 0.7;
    }
}

.notification {
    font-size: 22px;
    line-height: 14px;
    text-align: center;
    transition: color.3s;

    &:hover {
        color: grey;
        cursor: pointer;
    }
}

.main-content-page {
    padding: 0;
    background-color: #f7f8fa;
}

.dark .main-content-page {
    background-color: #212225;
}

.main-content {
    height: 100vh;
    width: 100vw;
}

.main-content-header {
    border-bottom: solid 1px var(--el-border-color);
    height: 55px;
    display: flex;
    align-items: center;
    box-sizing: border-box;

    .logo {
        height: 32px;
    }

    .user-info {
        display: flex;
        justify-content: flex-end;
        align-items: center;

        .el-avatar:hover {
            cursor: pointer;
        }

        .profile {
            text-align: right;
            margin-right: 20px;

            :first-child {
                font-size: 18px;
                font-weight: bold;
                line-height: 20px;
            }

            :last-child {
                font-size: 10px;
                color: grey;
            }
        }
    }
}

</style>
