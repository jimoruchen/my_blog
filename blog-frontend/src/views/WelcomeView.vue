<template>
    <div style="width: 100vw;height: 100vh;overflow: hidden;display: flex">
        <div style="flex: 1">
            <el-image style="width: 100%;height: 100%" fit="cover"
                      src="https://s2.loli.net/2025/04/03/1Sk4lByYNMTzPt7.jpg"/>
        </div>
        <div class="welcome-title">
            <div style="font-size: 30px;font-weight: bold">欢迎来到我的博客</div>
            <div style="margin-top: 10px">本平台旨在为您提供便捷的课程选择服务，确保每位学生都能根据自己的学术计划灵活安排学习进度。</div>
            <div style="margin-top: 5px">通过我们的系统，您可以轻松浏览、选择以及调整您的课程表，同时获取最新的课程相关信息与公告。</div>
        </div>
        <div class="right-card">
            <router-view v-slot="{ Component }">
                <transition name="el-fade-in-linear" mode="out-in">
                    <component :is="Component" style="height: 100%"/>
                </transition>
            </router-view>
        </div>
    </div>
</template>

<script setup>

</script>

<style scoped>
.right-card {
  width: 400px;
  z-index: 1;
  background-color: var(--el-bg-color);
}

.welcome-title {
    position: absolute;
    bottom: 30px;
    left: 30px;
    color: white;
    text-shadow: 0 0 10px black;
}
</style>
