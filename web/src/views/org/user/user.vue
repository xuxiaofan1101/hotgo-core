<template>
  <div>
    <n-card :bordered="false" class="proCard" title="用户管理">
      <n-tabs
        type="card"
        class="card-tabs"
        :value="defaultTab"
        animated
        @before-leave="handleBeforeLeave"
      >
        <n-tab-pane
          :name="item.key.toString()"
          :tab="item.label"
          v-for="item in dict.getOptionUnRef('roleTabs')"
          :key="item.key"
        >
          <List :type="defaultTab" />
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import List from './list.vue';
  import { useRouter } from 'vue-router';
  import { loadOptions } from './model';
  import { useDictStore } from '@/store/modules/dict';

  const dict = useDictStore();
  const router = useRouter();
  const defaultTab = ref('-1');

  onMounted(() => {
    if (router.currentRoute.value.query?.type) {
      defaultTab.value = router.currentRoute.value.query.type as string;
    }
    loadOptions();
  });

  function handleBeforeLeave(tabName: string): boolean | Promise<boolean> {
    defaultTab.value = tabName;
    return true;
  }
</script>
