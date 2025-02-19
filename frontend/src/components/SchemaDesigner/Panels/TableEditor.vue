<template>
  <div class="flex flex-col w-full h-full overflow-y-auto">
    <div class="py-2 w-full flex flex-row justify-between items-center">
      <div>
        <div
          v-if="!readonly"
          class="w-full flex justify-between items-center space-x-2"
        >
          <button
            class="flex flex-row justify-center items-center border px-3 py-1 leading-6 text-sm text-gray-700 rounded cursor-pointer hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="disableChangeTable"
            @click="handleAddColumn"
          >
            <heroicons-outline:plus class="w-4 h-auto mr-1 text-gray-400" />
            {{ $t("schema-editor.actions.add-column") }}
          </button>
          <button
            class="flex flex-row justify-center items-center border px-3 py-1 leading-6 text-sm text-gray-700 rounded cursor-pointer hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="disableChangeTable"
            @click="state.showSchemaTemplateDrawer = true"
          >
            <FeatureBadge feature="bb.feature.schema-template" />
            <heroicons-outline:plus class="w-4 h-auto mr-1 text-gray-400" />
            {{ $t("schema-editor.actions.add-from-template") }}
          </button>
        </div>
      </div>
      <div class="flex justify-end items-center">
        <NInput
          v-model:value="searchPattern"
          class="!w-48"
          :placeholder="$t('schema-editor.search-column')"
        >
          <template #prefix>
            <heroicons-outline:search class="w-4 h-auto text-gray-300" />
          </template>
        </NInput>
      </div>
    </div>

    <!-- column table -->
    <div
      id="table-editor-container"
      ref="tableEditorContainerRef"
      class="w-full h-auto grid auto-rows-auto border-y relative overflow-y-auto"
    >
      <!-- column table header -->
      <div
        class="sticky top-0 z-10 grid grid-cols-[6rem_minmax(0,_1.5fr)_repeat(3,_minmax(0,_0.8fr))_repeat(2,_80px)_minmax(0,_7rem)_20px] w-full text-sm leading-6 select-none bg-gray-50 text-gray-400"
        :class="shownColumnList.length > 0 && 'border-b'"
      >
        <span
          v-for="header in columnHeaderList"
          :key="header.key"
          class="table-header-item-container"
          >{{ header.label }}</span
        >
        <span></span>
      </div>
      <!-- column table body -->
      <div class="w-full">
        <div
          v-for="(column, index) in shownColumnList"
          :key="`${index}-${column.id}`"
          class="grid grid-cols-[6rem_minmax(0,_1.5fr)_repeat(3,_minmax(0,_0.8fr))_repeat(2,_80px)_minmax(0,_7rem)_20px] gr text-sm even:bg-gray-50"
          :class="[
            `column-${column.id}`,
            getColumnItemComputedClassList(column),
          ]"
        >
          <div class="table-body-item-container">
            <input
              v-model="column.name"
              :disabled="disableAlterColumn(column)"
              placeholder="column name"
              class="column-field-input column-name-input"
              type="text"
            />
          </div>
          <div
            class="table-body-item-container flex items-center gap-x-2 ml-3 text-sm"
          >
            {{ getColumnClassification(column)?.title ?? "N/A" }}
            <ClassificationLevelBadge
              :level-id="getColumnClassification(column)?.levelId"
              :classification-config="classificationConfig"
            />
            <div
              v-if="classificationConfig && !disableChangeTable"
              class="flex"
            >
              <button
                class="w-4 h-4 p-0.5 hover:bg-control-bg-hover rounded cursor-pointer"
                @click.prevent="column.classification = ''"
              >
                <heroicons-outline:x class="w-3 h-3" />
              </button>
              <button
                class="w-4 h-4 p-0.5 hover:bg-control-bg-hover rounded cursor-pointer"
                @click.prevent="state.pendingUpdateColumnIndex = index"
              >
                <heroicons-outline:pencil class="w-3 h-3" />
              </button>
            </div>
          </div>
          <div
            class="table-body-item-container flex flex-row justify-between items-center"
          >
            <input
              v-model="column.type"
              :disabled="
                disableAlterColumn(column) ||
                schemaTemplateColumnTypes.length > 0
              "
              placeholder="column type"
              class="column-field-input column-type-input !pr-8"
              type="text"
            />
            <NDropdown
              trigger="click"
              :disabled="disableAlterColumn(column)"
              :options="columnTypeOptions"
              @select="(dataType: string) => (column.type = dataType)"
            >
              <button class="absolute right-5">
                <heroicons-solid:chevron-up-down
                  class="w-4 h-auto text-gray-400"
                />
              </button>
            </NDropdown>
          </div>
          <div
            class="table-body-item-container flex flex-row justify-between items-center"
          >
            <input
              v-model="column.default"
              :disabled="disableAlterColumn(column)"
              :placeholder="column.default === undefined ? 'EMPTY' : 'NULL'"
              class="column-field-input !pr-8"
              type="text"
            />
            <NDropdown
              trigger="click"
              :disabled="disableAlterColumn(column)"
              :options="dataDefaultOptions"
              @select="(defaultString:string)=>handleColumnDefaultFieldChange(column, defaultString)"
            >
              <button class="absolute right-5">
                <heroicons-solid:chevron-up-down
                  class="w-4 h-auto text-gray-400"
                />
              </button>
            </NDropdown>
          </div>
          <div class="table-body-item-container">
            <input
              v-model="column.userComment"
              :disabled="disableAlterColumn(column)"
              placeholder="comment"
              class="column-field-input"
              type="text"
            />
          </div>
          <div
            class="table-body-item-container flex justify-start items-center"
          >
            <BBCheckbox
              class="ml-3"
              :value="!column.nullable"
              :disabled="
                disableAlterColumn(column) || isColumnPrimaryKey(column)
              "
              @toggle="(value) => (column.nullable = !value)"
            />
          </div>
          <div
            class="table-body-item-container flex justify-start items-center"
          >
            <BBCheckbox
              class="ml-3"
              :value="isColumnPrimaryKey(column)"
              :disabled="disableAlterColumn(column)"
              @toggle="(value) => setColumnPrimaryKey(column, value)"
            />
          </div>
          <div
            class="table-body-item-container foreign-key-field flex justify-start items-center"
          >
            <span
              v-if="checkColumnHasForeignKey(column)"
              class="column-field-text cursor-pointer !w-auto hover:opacity-80"
              @click="gotoForeignKeyReferencedTable(column)"
            >
              {{ getReferencedForeignKeyName(column) }}
            </span>
            <span v-else class="column-field-text italic text-gray-400 !w-auto"
              >EMPTY</span
            >
            <button
              v-if="!readonly"
              class="foreign-key-edit-button hidden cursor-pointer hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="disableAlterColumn(column)"
              @click="handleEditColumnForeignKey(column)"
            >
              <heroicons:pencil-square class="w-4 h-auto text-gray-400" />
            </button>
          </div>
          <div class="w-full flex justify-start items-center">
            <template v-if="!readonly">
              <n-tooltip v-if="!isDroppedColumn(column)" trigger="hover">
                <template #trigger>
                  <button
                    :disabled="disableChangeTable"
                    class="text-gray-500 cursor-pointer hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-60"
                    @click="handleDropColumn(column)"
                  >
                    <heroicons:trash class="w-4 h-auto" />
                  </button>
                </template>
                <span>{{ $t("schema-editor.actions.drop-column") }}</span>
              </n-tooltip>
              <n-tooltip v-else trigger="hover">
                <template #trigger>
                  <button
                    class="text-gray-500 cursor-pointer hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="disableChangeTable"
                    @click="handleRestoreColumn(column)"
                  >
                    <heroicons:arrow-uturn-left class="w-4 h-auto" />
                  </button>
                </template>
                <span>{{ $t("schema-editor.actions.restore") }}</span>
              </n-tooltip>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>

  <EditColumnForeignKeyModal
    v-if="state.showEditColumnForeignKeyModal && editForeignKeyColumn"
    :schema-id="schema.id"
    :table-id="table.id"
    :column-id="editForeignKeyColumn.id"
    @close="state.showEditColumnForeignKeyModal = false"
  />

  <Drawer
    :show="state.showSchemaTemplateDrawer"
    @close="state.showSchemaTemplateDrawer = false"
  >
    <DrawerContent :title="$t('schema-template.field-template.self')">
      <div class="w-[calc(100vw-36rem)] min-w-[64rem] max-w-[calc(100vw-8rem)]">
        <FieldTemplates :engine="engine" @apply="handleApplyColumnTemplate" />
      </div>
    </DrawerContent>
  </Drawer>
  <FeatureModal
    feature="bb.feature.schema-template"
    :open="state.showFeatureModal"
    @cancel="state.showFeatureModal = false"
  />

  <SelectClassificationDrawer
    v-if="classificationConfig"
    :show="state.pendingUpdateColumnIndex >= 0"
    :classification-config="classificationConfig"
    @dismiss="state.pendingUpdateColumnIndex = -1"
    @select="onClassificationSelect"
  />
</template>

<script lang="ts" setup>
import { isUndefined, flatten } from "lodash-es";
import { NDropdown } from "naive-ui";
import scrollIntoView from "scroll-into-view-if-needed";
import { computed, nextTick, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import { BBCheckbox } from "@/bbkit";
import { Drawer, DrawerContent } from "@/components/v2";
import {
  hasFeature,
  generateUniqueTabId,
  useSettingV1Store,
} from "@/store/modules";
import { ColumnMetadata } from "@/types/proto/store/database";
import { Engine } from "@/types/proto/v1/common";
import { SchemaTemplateSetting_FieldTemplate } from "@/types/proto/v1/setting_service";
import {
  Column,
  Table,
  Schema,
  convertColumnMetadataToColumn,
  ForeignKey,
} from "@/types/schemaEditor/atomType";
import { getDataTypeSuggestionList } from "@/utils";
import FieldTemplates from "@/views/SchemaTemplate/FieldTemplates.vue";
import EditColumnForeignKeyModal from "../Modals/EditColumnForeignKeyModal.vue";
import {
  SchemaDesignerTabType,
  TableTabContext,
  useSchemaDesignerContext,
} from "../common";
import { isColumnChanged } from "../utils/column";

interface LocalState {
  isFetchingDDL: boolean;
  statement: string;
  showEditColumnForeignKeyModal: boolean;
  showSchemaTemplateDrawer: boolean;
  showFeatureModal: boolean;
  pendingUpdateColumnIndex: number;
}

const { t } = useI18n();
const { readonly, engine, project, editableSchemas, getCurrentTab, addTab } =
  useSchemaDesignerContext();
useSchemaDesignerContext();
const settingStore = useSettingV1Store();
const currentTab = computed(() => getCurrentTab() as TableTabContext);
const state = reactive<LocalState>({
  isFetchingDDL: false,
  statement: "",
  showEditColumnForeignKeyModal: false,
  showSchemaTemplateDrawer: false,
  showFeatureModal: false,
  pendingUpdateColumnIndex: -1,
});

const schema = computed(() => {
  return editableSchemas.value.find(
    (schema) => schema.id === currentTab.value.schemaId
  ) as Schema;
});
const table = computed(
  () =>
    schema.value.tableList.find(
      (table) => table.id === currentTab.value.tableId
    ) as Table
);
const foreignKeyList = computed(() => {
  return schema.value.foreignKeyList.filter(
    (pk) => pk.tableId === currentTab.value.tableId
  ) as ForeignKey[];
});

const searchPattern = ref("");
const tableEditorContainerRef = ref<HTMLDivElement>();
const editForeignKeyColumn = ref<Column>();

const shownColumnList = computed(() => {
  return table.value.columnList.filter((column) =>
    column.name.includes(searchPattern.value.trim())
  );
});
const isDroppedSchema = computed(() => {
  return schema.value.status === "dropped";
});

const isDroppedTable = computed(() => {
  return table.value.status === "dropped";
});

const columnHeaderList = computed(() => {
  return [
    {
      key: "name",
      label: t("schema-editor.column.name"),
    },
    {
      key: "classification",
      label: t("schema-editor.column.classification"),
    },
    {
      key: "type",
      label: t("schema-editor.column.type"),
    },
    {
      key: "default",
      label: t("schema-editor.column.default"),
    },
    {
      key: "comment",
      label: t("schema-editor.column.comment"),
    },
    {
      key: "nullable",
      label: t("schema-editor.column.not-null"),
    },
    {
      key: "primary",
      label: t("schema-editor.column.primary"),
    },
    {
      key: "foreign_key",
      label: t("schema-editor.column.foreign-key"),
    },
  ];
});

const schemaTemplateColumnTypes = computed(() => {
  const setting = settingStore.getSettingByName("bb.workspace.schema-template");
  const columnTypes = setting?.value?.schemaTemplateSettingValue?.columnTypes;
  if (columnTypes && columnTypes.length > 0) {
    const columnType = columnTypes.find(
      (columnType) => columnType.engine === engine.value
    );
    if (columnType && columnType.enabled) {
      return columnType.types;
    }
  }
  return [];
});

const columnTypeOptions = computed(() => {
  if (schemaTemplateColumnTypes.value.length > 0) {
    return schemaTemplateColumnTypes.value.map((columnType) => {
      return {
        label: columnType,
        key: columnType,
      };
    });
  }

  return getDataTypeSuggestionList(engine.value).map((dataType) => {
    return {
      label: dataType,
      key: dataType,
    };
  });
});

const getColumnItemComputedClassList = (column: Column) => {
  if (column.status === "dropped") {
    return ["text-red-700", "cursor-not-allowed", "!bg-red-50", "opacity-70"];
  } else if (column.status === "created") {
    return ["text-green-700", "!bg-green-50"];
  } else if (
    isColumnChanged(
      currentTab.value.schemaId,
      currentTab.value.tableId,
      column.id
    )
  ) {
    return ["text-yellow-700", "!bg-yellow-50"];
  }
  return [];
};

const dataDefaultOptions = [
  {
    label: "NULL",
    key: "NULL",
  },
  {
    label: "EMPTY",
    key: "EMPTY",
  },
];

const isColumnPrimaryKey = (column: Column): boolean => {
  return table.value.primaryKey.columnIdList.includes(column.id);
};

const checkColumnHasForeignKey = (column: Column): boolean => {
  const columnIdList = flatten(
    foreignKeyList.value.map((fk) => fk.columnIdList)
  );
  return columnIdList.includes(column.id);
};

const getReferencedForeignKeyName = (column: Column) => {
  if (!checkColumnHasForeignKey(column)) {
    return;
  }
  const fk = foreignKeyList.value.find(
    (fk) =>
      fk.columnIdList.find((columnId) => columnId === column.id) !== undefined
  );
  const index = fk?.columnIdList.findIndex(
    (columnId) => columnId === column.id
  );

  if (isUndefined(fk) || isUndefined(index) || index < 0) {
    return;
  }
  const referencedSchema = editableSchemas.value.find(
    (schema) => schema.id === fk.referencedSchemaId
  );
  const referencedTable = referencedSchema?.tableList.find(
    (table) => table.id === fk.referencedTableId
  );
  if (!referencedTable) {
    return;
  }
  const referColumn = referencedTable.columnList.find(
    (column) => column.id === fk.referencedColumnIdList[index]
  );
  if (engine.value === Engine.MYSQL) {
    return `${referencedTable.name}(${referColumn?.name})`;
  } else {
    return `${referencedSchema?.name}.${referencedTable.name}(${referColumn?.name})`;
  }
};

const isDroppedColumn = (column: Column): boolean => {
  return column.status === "dropped";
};

const disableChangeTable = computed(() => {
  return readonly.value || isDroppedSchema.value || isDroppedTable.value;
});

const disableAlterColumn = (column: Column): boolean => {
  return (
    readonly.value ||
    isDroppedSchema.value ||
    isDroppedTable.value ||
    isDroppedColumn(column)
  );
};

const setColumnPrimaryKey = (column: Column, isPrimaryKey: boolean) => {
  if (isPrimaryKey) {
    column.nullable = false;
    table.value.primaryKey.columnIdList.push(column.id);
  } else {
    table.value.primaryKey.columnIdList =
      table.value.primaryKey.columnIdList.filter(
        (columnId) => columnId !== column.id
      );
  }
};

const handleAddColumn = () => {
  const column = convertColumnMetadataToColumn(ColumnMetadata.fromPartial({}));
  column.status = "created";
  table.value.columnList.push(column);
  nextTick(() => {
    (
      tableEditorContainerRef.value?.querySelector(
        `.column-${column.id} .column-name-input`
      ) as HTMLInputElement
    )?.focus();
  });
};

const handleApplyColumnTemplate = (
  template: SchemaTemplateSetting_FieldTemplate
) => {
  if (!hasFeature("bb.feature.schema-template")) {
    state.showFeatureModal = true;
    return;
  }
  if (template.engine !== engine.value || !template.column) {
    return;
  }
  const column = convertColumnMetadataToColumn(template.column);
  column.status = "created";
  table.value.columnList.push(column);
  state.showSchemaTemplateDrawer = false;
};

const handleColumnDefaultFieldChange = (
  column: Column,
  defaultString: string
) => {
  if (defaultString === "NULL") {
    column.default = "NULL";
  } else if (defaultString === "EMPTY") {
    column.default = undefined;
  }
};

const gotoForeignKeyReferencedTable = (column: Column) => {
  if (!checkColumnHasForeignKey(column)) {
    return;
  }
  const fk = foreignKeyList.value.find(
    (fk) =>
      fk.columnIdList.find((columnId) => columnId === column.id) !== undefined
  );
  const index = fk?.columnIdList.findIndex(
    (columnId) => columnId === column.id
  );
  if (isUndefined(fk) || isUndefined(index) || index < 0) {
    return;
  }

  const referencedSchema = editableSchemas.value.find(
    (schema) => schema.id === fk.referencedSchemaId
  );
  const referencedTable = referencedSchema?.tableList.find(
    (table) => table.id === fk.referencedTableId
  );
  if (!referencedTable) {
    return;
  }
  const referColumn = referencedTable?.columnList.find(
    (column) => column.id === fk.referencedColumnIdList[index]
  );
  if (!referencedSchema || !referencedTable || !referColumn) {
    return;
  }

  addTab({
    id: generateUniqueTabId(),
    type: SchemaDesignerTabType.TabForTable,
    schemaId: referencedSchema.id,
    tableId: referencedTable.id,
  });

  nextTick(() => {
    const container = document.querySelector("#table-editor-container");
    const input = container?.querySelector(
      `.column-${referColumn.id} .column-name-input`
    ) as HTMLInputElement | undefined;
    if (input) {
      input.focus();
      scrollIntoView(input);
    }
  });
};

const handleEditColumnForeignKey = (column: Column) => {
  editForeignKeyColumn.value = column;
  state.showEditColumnForeignKeyModal = true;
};

const handleDropColumn = (column: Column) => {
  if (column.status === "created") {
    table.value.columnList = table.value.columnList.filter(
      (item) => item !== column
    );
    table.value.primaryKey.columnIdList =
      table.value.primaryKey.columnIdList.filter(
        (columnId) => columnId !== column.id
      );

    const foreignKeyList = schema.value.foreignKeyList.filter(
      (fk) => fk.tableId === currentTab.value.tableId
    );
    for (const foreignKey of foreignKeyList) {
      const columnRefIndex = foreignKey.columnIdList.findIndex(
        (columnId) => columnId === column.id
      );
      if (columnRefIndex > -1) {
        foreignKey.columnIdList.splice(columnRefIndex, 1);
        foreignKey.referencedColumnIdList.splice(columnRefIndex, 1);
      }
    }
  } else {
    column.status = "dropped";
  }
};

const handleRestoreColumn = (column: Column) => {
  if (column.status === "created") {
    return;
  }

  column.status = "normal";
};

const classificationConfig = computed(() => {
  if (!project.value || !project.value.dataClassificationConfigId) {
    return;
  }
  return settingStore.getProjectClassification(
    project.value.dataClassificationConfigId
  );
});

const getColumnClassification = (column: Column) => {
  if (!classificationConfig.value) {
    return;
  }
  const { classification } = column;
  if (!classification) {
    return;
  }
  return classificationConfig.value.classification[classification];
};

const onClassificationSelect = (classificationId: string) => {
  if (!table.value.columnList[state.pendingUpdateColumnIndex]) {
    return;
  }
  table.value.columnList[state.pendingUpdateColumnIndex].classification =
    classificationId;
  state.pendingUpdateColumnIndex = -1;
};
</script>

<style scoped>
.table-header-item-container {
  @apply py-2 px-3;
}
.table-body-item-container {
  @apply w-full h-10 box-border p-px pr-2 relative;
}
.column-field-input {
  @apply w-full pr-1 box-border border-transparent truncate select-none rounded bg-transparent text-sm placeholder:italic placeholder:text-gray-400 focus:bg-white focus:text-black;
}
.column-field-text {
  @apply w-full pl-3 pr-1 box-border border-transparent truncate select-none rounded bg-transparent text-sm placeholder:italic placeholder:text-gray-400 focus:bg-white focus:text-black;
}
.foreign-key-field:hover .foreign-key-edit-button {
  @apply block;
}
</style>
