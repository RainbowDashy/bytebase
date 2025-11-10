<template>
  <div
    v-for="(error, i) in syntaxErrors"
    :key="i"
    class="text-sm grid gap-1 pl-8"
    style="grid-template-columns: auto 1fr"
  >
    <template v-if="error.position">
      <div>{{ $t("common.position") }}:</div>
      <div>
        {{
          $t("sql-editor.syntax-error-position", {
            line: error.position.line,
            column: error.position.column,
          })
        }}
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { SQLResultSetV1 } from "@/types";
import type { QueryResult_SyntaxErrorDetail } from "@/types/proto-es/v1/sql_service_pb";

const props = defineProps<{
  resultSet: SQLResultSetV1;
}>();

const syntaxErrors = computed(() => {
  const errors: QueryResult_SyntaxErrorDetail[] = [];
  props.resultSet.results.forEach((result) => {
    if (result.detailedError?.case === "syntaxError") {
      errors.push(result.detailedError.value);
    }
  });
  return errors;
});
</script>
