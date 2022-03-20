import { defineComponent, nextTick, ref } from "vue";
import {
  FormRules,
  NButton,
  NCard,
  NPageHeader,
  NSpin,
  NText,
  NUpload,
  NUploadDragger,
  useMessage,
  NGridItem,
  UploadFileInfo,
  UploadInst,
  NIcon,
} from "naive-ui";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import { css } from "@linaria/core";
import { Upload } from "@vicons/fa";

import { Mode } from "../../states/common";
import { padding } from "../../constants/style";
import useImageState, { imageList } from "../../states/image";
import ExTable from "../../components/ExTable";
import { newRequireRule } from "../../components/ExConfigEditor";
import ExForm, { FormItemTypes, FormItem } from "../../components/ExForm";
import { showError } from "../../helpers/util";
import { getAPIUrl } from "../../helpers/request";
import { IMAGES } from "../../constants/url";

const addButtonClass = css`
  width: 100%;
  margin-top: ${2 * padding}px;
`;

function getColumns(): TableColumn[] {
  return [
    {
      title: "名称",
      key: "name",
    },
    {
      title: "bucket",
      key: "bucket",
    },
    {
      title: "类型",
      key: "type",
    },
    {
      title: "宽度",
      key: "width",
    },
    {
      title: "高度",
      key: "height",
    },
    {
      title: "标签",
      key: "tags",
    },
    {
      title: "描述",
      key: "description",
    },
    {
      title: "更新时间",
      key: "updatedAt",
    },
  ];
}

function getFormItems(params: {
  bucket?: string;
  name?: string;
  tags?: string;
  description?: string;
}): FormItem[] {
  const defaultBucketValue: string[] = [];
  if (params.bucket) {
    defaultBucketValue.push(params.bucket);
  }
  return [
    {
      type: FormItemTypes.BucketSelect,
      name: "Bucket:",
      key: "bucket",
      defaultValue: params.bucket,
      placeholder: "请先选择图片bucket",
    },
    {
      name: "名称：",
      key: "name",
      defaultValue: params.name,
      placeholder: "请输入图片名称(可选)",
    },
    {
      name: "标签：",
      key: "tags",
      defaultValue: params.tags,
      placeholder: "请输入图片标签，以,分隔",
    },
    {
      name: "描述：",
      key: "description",
      defaultValue: params.description,
      span: 16,
      type: FormItemTypes.TextArea,
      placeholder: "请输入该图片的描述",
    },
  ];
}

export default defineComponent({
  name: "ImageList",
  setup() {
    const message = useMessage();
    const { images } = useImageState();

    const mode = ref(Mode.List);
    const updatedID = ref(0);
    const currentImage = ref<Record<string, unknown>>({});
    const processing = ref(false);
    const uploadRef = ref<UploadInst | null>(null);

    const toggle = (value: Mode) => {
      if (value === Mode.List) {
        currentImage.value = {
          bucket: currentImage.value.bucket,
        };
        updatedID.value = 0;
      }
      mode.value = value;
    };

    let uploadCount = 0;
    const uploadAction = ref("");
    const onSubmit = async (data: Record<string, unknown>) => {
      if (uploadCount <= 0) {
        showError(message, new Error("请先上传图片"));
        return;
      }
      const arr: string[] = [];
      ["bucket", "name", "tags", "description"].forEach((key) => {
        if (!data[key]) {
          return;
        }
        arr.push(`${key}=${data[key]}`);
      });
      uploadAction.value = `${getAPIUrl(IMAGES)}?${arr.join("&")}`;
      await nextTick();
      uploadRef.value?.submit();
      processing.value = true;
    };
    const onFinish = ({
      file,
      event,
    }: {
      file: UploadFileInfo;
      event?: ProgressEvent;
    }) => {
      message.success((event?.target as XMLHttpRequest).response);
      processing.value = false;
      toggle(Mode.List);
      return file;
    };

    const onChange = (data: { fileList: UploadFileInfo[] }) => {
      uploadCount = data.fileList.length;
    };

    const fetch = (params: {
      bucket: string;
      tag: string;
      limit: number;
      offset: number;
    }) => {
      if (!params.bucket) {
        return;
      }
      currentImage.value.bucket = params.bucket;
      return imageList(params);
    };

    return {
      uploadAction,
      upload: uploadRef,
      images,
      mode,
      processing,
      updatedID,
      currentImage,
      fetch,
      toggle,
      onSubmit,
      onChange,
      onFinish,
    };
  },
  render() {
    const {
      mode,
      images,
      processing,
      toggle,
      fetch,
      currentImage,
      updatedID,
      uploadAction,
      onSubmit,
      onChange,
      onFinish,
    } = this;
    if (mode == Mode.List) {
      const columns = getColumns();
      const filters = [
        {
          type: FormItemTypes.BucketSelect,
          defaultValue: currentImage.bucket,
          placeholder: "请先选择图片bucket",
          key: "bucket",
        },
      ];
      return (
        <NCard title={"图片列表"}>
          <ExTable
            columns={columns}
            data={images}
            fetch={fetch}
            filters={filters}
          />
          <NButton
            size="large"
            class={addButtonClass}
            onClick={() => {
              toggle(Mode.Add);
            }}
          >
            上传图片
          </NButton>
        </NCard>
      );
    }
    const rules: FormRules = {
      bucket: newRequireRule("bucket不能为空"),
      description: newRequireRule("描述不能为空"),
    };
    const uploadComponent = (
      <NUpload
        ref="upload"
        defaultUpload={false}
        onChange={onChange}
        onFinish={onFinish}
        action={uploadAction}
      >
        <NUploadDragger>
          <NText>
            <NIcon class="mright5" size={16}>
              <Upload />
            </NIcon>
            点击或拖动图片至此区域上传
          </NText>
        </NUploadDragger>
      </NUpload>
    );
    const formItems = getFormItems(currentImage);
    formItems.push({
      name: "上传图片：",
      key: "file",
      component: uploadComponent,
    });
    return (
      <NSpin show={processing}>
        <NCard>
          <NPageHeader
            title={"更新/添加图片"}
            onBack={() => {
              toggle(Mode.List);
            }}
          >
            <ExForm
              formItems={formItems}
              submitText={updatedID !== 0 ? "更新" : "添加"}
              onSubmit={onSubmit}
              rules={rules}
            >
              <NGridItem span={24}></NGridItem>
            </ExForm>
          </NPageHeader>
        </NCard>
      </NSpin>
    );
  },
});
