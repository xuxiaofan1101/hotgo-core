<template>
  <div class="console">
    <!--数据卡片-->
    <n-grid cols="1 s:2 m:3 l:4 xl:4 2xl:4" responsive="screen" :x-gap="12" :y-gap="8">
      <n-grid-item>
        <NCard
          title="卡板量"
          :segmented="{ content: true, footer: true }"
          size="small"
          :bordered="false"
        >
          <template #header-extra>
            <n-tag type="success">日</n-tag>
          </template>
          <div class="py-1 px-1 flex justify-between">
            <n-skeleton v-if="loading" :width="100" size="medium" />
            <CountTo v-else :startVal="1" :endVal="visits.dayVisits" class="text-3xl" />
          </div>
          <div class="py-1 px-1 flex justify-between">
            <div class="text-sn">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                日同比
                <CountTo :startVal="1" suffix="%" :endVal="visits.rise" />
                <n-icon size="12" color="#00ff6f">
                  <CaretUpOutlined />
                </n-icon>
              </template>
            </div>
            <div class="text-sn">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                周同比
                <CountTo :startVal="1" suffix="%" :endVal="visits.decline" />
                <n-icon size="12" color="#ffde66">
                  <CaretDownOutlined />
                </n-icon>
              </template>
            </div>
          </div>
          <template #footer>
            <div class="flex justify-between">
              <n-skeleton v-if="loading" text :repeat="2" />
              <template v-else>
                <div class="text-sn"> 总卡板量： </div>
                <div class="text-sn">
                  <CountTo :startVal="1" :endVal="visits.amount" />
                </div>
              </template>
            </div>
          </template>
        </NCard>
      </n-grid-item>
      <n-grid-item>
        <NCard
          title="激活卡板"
          :segmented="{ content: true, footer: true }"
          size="small"
          :bordered="false"
        >
          <template #header-extra>
            <n-tag type="info">周</n-tag>
          </template>
          <div class="py-1 px-1 flex justify-between">
            <n-skeleton v-if="loading" :width="100" size="medium" />
            <CountTo
              v-else
              prefix="￥"
              :startVal="1"
              :endVal="saleroom.weekSaleroom"
              class="text-3xl"
            />
          </div>
          <div class="py-2 px-2 flex justify-between">
            <div class="text-sn flex-1">
              <n-progress
                type="line"
                :percentage="saleroom.degree"
                :indicator-placement="'inside'"
                processing
              />
            </div>
          </div>
          <template #footer>
            <div class="flex justify-between">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                <div class="text-sn"> 总激活卡板： </div>
                <div class="text-sn">
                  <CountTo :startVal="1" :endVal="saleroom.amount" />
                  <!-- prefix="￥"-->
                </div>
              </template>
            </div>
          </template>
        </NCard>
      </n-grid-item>
      <n-grid-item>
        <NCard
          title="代理商"
          :segmented="{ content: true, footer: true }"
          size="small"
          :bordered="false"
        >
          <template #header-extra>
            <n-tag type="warning">周</n-tag>
          </template>
          <div class="py-1 px-1 flex justify-between">
            <n-skeleton v-if="loading" :width="100" size="medium" />
            <CountTo v-else :startVal="1" :endVal="orderLarge.weekLarge" class="text-3xl" />
          </div>
          <div class="py-1 px-1 flex justify-between">
            <div class="text-sn">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                日同比
                <CountTo :startVal="1" suffix="%" :endVal="orderLarge.rise" />
                <n-icon size="12" color="#00ff6f">
                  <CaretUpOutlined />
                </n-icon>
              </template>
            </div>
            <div class="text-sn">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                周同比
                <CountTo :startVal="1" suffix="%" :endVal="orderLarge.rise" />
                <n-icon size="12" color="#ffde66">
                  <CaretDownOutlined />
                </n-icon>
              </template>
            </div>
          </div>
          <template #footer>
            <div class="flex justify-between">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                <div class="text-sn"> 总代理商量： </div>
                <div class="text-sn">
                  <CountTo :startVal="1" suffix="%" :endVal="orderLarge.amount" />
                </div>
              </template>
            </div>
          </template>
        </NCard>
      </n-grid-item>
      <n-grid-item>
        <NCard
          title="提现佣金"
          :segmented="{ content: true, footer: true }"
          size="small"
          :bordered="false"
        >
          <template #header-extra>
            <n-tag type="error">月</n-tag>
          </template>
          <div class="py-1 px-1 flex justify-between">
            <n-skeleton v-if="loading" :width="100" size="medium" />
            <CountTo v-else prefix="￥" :startVal="1" :endVal="volume.weekLarge" class="text-3xl" />
          </div>
          <div class="py-1 px-1 flex justify-between">
            <div class="text-sn">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                月同比
                <CountTo :startVal="1" suffix="%" :endVal="volume.rise" />
                <n-icon size="12" color="#00ff6f">
                  <CaretUpOutlined />
                </n-icon>
              </template>
            </div>
            <div class="text-sn">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                月同比
                <CountTo :startVal="1" suffix="%" :endVal="volume.decline" />
                <n-icon size="12" color="#ffde66">
                  <CaretDownOutlined />
                </n-icon>
              </template>
            </div>
          </div>
          <template #footer>
            <div class="flex justify-between">
              <n-skeleton v-if="loading" :width="100" size="medium" />
              <template v-else>
                <div class="text-sn"> 总提现额： </div>
                <div class="text-sn">
                  <CountTo prefix="￥" :startVal="1" :endVal="volume.amount" />
                </div>
              </template>
            </div>
          </template>
        </NCard>
      </n-grid-item>
    </n-grid>

    <!--导航卡片-->
    <div class="mt-4">
      <n-grid cols="1 s:2 m:3 l:8 xl:8 2xl:8" responsive="screen" :x-gap="16" :y-gap="8">
        <n-grid-item v-for="(item, index) in iconList" :key="index" @click="item.eventObject || {}">
          <NCard
            content-style="padding-top: 0;"
            size="small"
            :bordered="false"
            v-on="item.eventObject || {}"
          >
            <template #footer>
              <n-skeleton v-if="loading" size="medium" />
              <div class="cursor-pointer" v-else>
                <p class="flex justify-center">
                  <span>
                    <n-icon :size="item.size" class="flex-1" :color="item.color">
                      <component :is="item.icon" />
                    </n-icon>
                  </span>
                </p>
                <p class="flex justify-center"
                  ><span>{{ item.title }}</span></p
                >
              </div>
            </template>
          </NCard>
        </n-grid-item>
      </n-grid>
    </div>

    <!--访问量 | 流量趋势-->
    <VisiTab />
  </div>
</template>
<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { getConsoleInfo } from '@/api/dashboard/console';
  import VisiTab from './components/VisiTab.vue';
  import { CountTo } from '@/components/CountTo';
  import {
    CaretUpOutlined,
    CaretDownOutlined,
    UsergroupAddOutlined,
    BarChartOutlined,
    ShoppingCartOutlined,
    AccountBookOutlined,
    CreditCardOutlined,
    MailOutlined,
    TagsOutlined,
    SettingOutlined,
  } from '@vicons/antd';

  const loading = ref(true);
  const visits = ref<any>({});
  const saleroom = ref<any>({});
  const orderLarge = ref<any>({});
  const volume = ref({});
  const router = useRouter();

  // 图标列表
  const iconList = [
    {
      icon: UsergroupAddOutlined,
      size: '32',
      title: '用户',
      color: '#69c0ff',
      eventObject: {
        click: () => router.push({ name: 'user' }),
      },
    },
    {
      icon: BarChartOutlined,
      size: '32',
      title: '分析',
      color: '#69c0ff',
      eventObject: {
        click: () => {},
      },
    },
    {
      icon: ShoppingCartOutlined,
      size: '32',
      title: '商品',
      color: '#ff9c6e',
      eventObject: {
        click: () => {},
      },
    },
    {
      icon: AccountBookOutlined,
      size: '32',
      title: '订单',
      color: '#b37feb',
      eventObject: {
        click: () => {},
      },
    },
    {
      icon: CreditCardOutlined,
      size: '32',
      title: '票据',
      color: '#ffd666',
      eventObject: {
        click: () => {},
      },
    },
    {
      icon: MailOutlined,
      size: '32',
      title: '消息',
      color: '#5cdbd3',
      eventObject: {
        click: () => {},
      },
    },
    {
      icon: TagsOutlined,
      size: '32',
      title: '标签',
      color: '#ff85c0',
      eventObject: {
        click: () => {},
      },
    },
    {
      icon: SettingOutlined,
      size: '32',
      title: '配置',
      color: '#ffc069',
      eventObject: {
        click: () => {},
      },
    },
  ];

  onMounted(async () => {
    const data = await getConsoleInfo();
    visits.value = data.visits;
    saleroom.value = data.saleroom;
    orderLarge.value = data.orderLarge;
    volume.value = data.volume;
    loading.value = false;
  });
</script>

<style lang="less" scoped></style>
