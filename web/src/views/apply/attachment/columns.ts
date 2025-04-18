import { h, ref } from 'vue';
import { NAvatar, NImage } from 'naive-ui';
import { getFileExt } from '@/utils/urlUtils';
import { FormSchema } from '@/components/Form';
import { fallbackSrc } from '@/utils/hotgo';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { useDictStore } from '@/store/modules/dict';
import { renderOptionTag } from '@/utils';

const dict = useDictStore();

export const schemas = ref<FormSchema[]>([
  {
    field: 'drive',
    component: 'NSelect',
    label: '上传驱动',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择上传驱动',
      options: dict.getOption('config_upload_drive'),
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
  {
    field: 'kind',
    component: 'NSelect',
    label: '上传类型',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择上传类型',
      options: dict.getOption('AttachmentKindOption'),
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
  {
    field: 'name',
    component: 'NInput',
    label: '文件名称',
    labelMessage: '支持模糊查询',
    componentProps: {
      placeholder: '请输入文件名称',
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
    rules: [{ message: '请输入文件名称', trigger: ['blur'] }],
  },
  {
    field: 'updatedAt',
    component: 'NDatePicker',
    label: '上传时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
  {
    field: 'member_id',
    component: 'NInput',
    label: '用户ID',
    componentProps: {
      placeholder: '请输入用户ID',
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
    rules: [{ message: '请输入用户ID', trigger: ['blur'] }],
  },
]);

export const columns = [
  {
    title: '附件ID',
    key: 'id',
    width: 80,
  },
  {
    title: '应用',
    key: 'appId',
    width: 100,
  },
  {
    title: '驱动',
    key: 'drive',
    render(row) {
      return row.drive;
    },
    width: 100,
  },
  {
    title: '上传类型',
    key: 'kind',
    render(row) {
      return renderOptionTag('AttachmentKindOption', row.kind);
    },
    width: 120,
  },
  {
    title: '文件名称',
    key: 'name',
    width: 150,
  },
  {
    title: '文件',
    key: 'fileUrl',
    width: 80,
    render(row) {
      if (row.fileUrl === '') {
        return ``;
      }
      if (row.kind !== 'image') {
        return h(
          NAvatar,
          {
            width: '40px',
            height: '40px',
            'max-width': '100%',
            'max-height': '100%',
          },
          {
            default: () => getFileExt(row.fileUrl),
          }
        );
      }
      return h(NImage, {
        width: 40,
        height: 40,
        src: row.fileUrl,
        fallbackSrc: fallbackSrc(),
        style: {
          width: '40px',
          height: '40px',
          'max-width': '100%',
          'max-height': '100%',
        },
      });
    },
  },
  {
    title: '文件大小',
    key: 'sizeFormat',
    width: 100,
  },
  {
    title: '扩展名',
    key: 'ext',
    width: 80,
  },
  {
    title: '扩展类型',
    key: 'mimeType',
    width: 200,
  },
  {
    title: '上传时间',
    key: 'updatedAt',
    width: 180,
  },
];

export function loadOptions() {
  dict.loadOptions(['AttachmentKindOption', 'config_upload_drive']);
}
