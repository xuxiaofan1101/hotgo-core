<script lang="ts" setup>
  import { ref, onMounted, onBeforeUnmount } from 'vue';
  import { useRouter, onBeforeRouteLeave } from 'vue-router';

  const router = useRouter();
  const countdown = ref(5);
  let intervalId: ReturnType<typeof setInterval>;

  function goHome() {
    router.push('/');
  }

  const startCountdown = () => {
    intervalId = setInterval(() => {
      if (countdown.value > 0) {
        countdown.value--;
      } else {
        router.push('/');
        clearInterval(intervalId);
      }
    }, 1000);
  };

  onMounted(startCountdown);
  onBeforeUnmount(() => clearInterval(intervalId));
  onBeforeRouteLeave(() => clearInterval(intervalId));
</script>

<template>
  <div class="flex flex-col justify-center page-container">
    <div class="text-center">
      <img src="~@/assets/images/exception/404.svg" alt="" />
    </div>
    <div class="text-center">
      <h1 class="text-base text-gray-500">
        抱歉，页面不存在！<br />
        {{ countdown }}秒后自动返回首页
      </h1>
      <n-button type="info" @click="goHome">回到首页</n-button>
    </div>
  </div>
</template>

<style lang="less" scoped>
  .page-container {
    width: 100%;
    border-radius: 4px;
    padding: 50px 0;
    height: 100vh;

    .text-center {
      h1 {
        color: #666;
        padding: 20px 0;
      }
    }

    img {
      width: 350px;
      margin: 0 auto;
    }
  }
</style>
