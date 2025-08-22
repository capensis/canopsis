<template>
  <div>
    <v-layout v-if="!hasStructure" class="gap-4" column>
      <v-flex align-self-center>
        <span class="grey--text">{{ $t('common.noDataAvailable') }}</span>
      </v-flex>
      <v-flex align-self-center>
        <v-btn
          color="primary"
          @click="$emit('import')"
        >
          {{ $t('common.importData') }}
        </v-btn>
      </v-flex>
    </v-layout>
    <v-layout v-else-if="isEmptyColumns" column>
      <v-flex align-self-center>
        <v-subheader>{{ $t('externalData.tableEmptyColumns') }}</v-subheader>
      </v-flex>
    </v-layout>
    <c-advanced-data-table
      v-else
      ref="tableElement"
      :loading="pending"
      :options="options"
      :items="records"
      :headers="sortedHeadersWithWidth"
      :expand="expandable"
      :select-all="selectable"
      :total-items="totalItems"
      :search="withToolbar"
      :dense="isMediumDense"
      :ultra-dense="isSmallDense"
      :table-class="tableClass"
      :style="tableStyle"
      advanced-pagination
      @update:options="$emit('update:options', $event)"
    >
      <template v-if="withToolbar" #toolbar="{ selected }">
        <v-flex>
          <v-layout class="gap-1" justify-end>
            <v-btn
              color="blue"
              fab
              small
              @click="$emit('import')"
            >
              <v-icon color="white">
                upload_file
              </v-icon>
            </v-btn>
            <v-btn
              color="primary"
              fab
              small
              @click="$emit('add')"
            >
              <v-icon color="white">
                add
              </v-icon>
            </v-btn>
          </v-layout>
        </v-flex>
        <v-flex class="py-2" xs12>
          <v-layout align-center>
            <v-flex xs4 align-self-center>
              <c-density-btn-toggle
                v-if="densable"
                :value="dense"
                @change="$emit('update:dense', $event)"
              />
            </v-flex>
            <v-flex xs4 align-self-center>
              <c-pagination
                :page="options.page"
                :limit="options.itemsPerPage"
                :total="totalItems"
                type="top"
                @input="updatePage"
              />
            </v-flex>
            <v-flex class="text-right" xs4 align-self-center>
              <v-layout>
                <c-grid-btns
                  v-if="resizable || draggable"
                  :editing="isColumnsEditing"
                  @toggle="toggleColumnEditingMode"
                  @reset="resetColumnsSettings"
                />
                <c-action-btn
                  v-if="exportable"
                  :loading="downloading"
                  :tooltip="$t('settings.exportAsCsv')"
                  icon="cloud_download"
                  @click="$emit('export', selected)"
                />
              </v-layout>
            </v-flex>
          </v-layout>
        </v-flex>
      </template>
      <template #mass-actions="{ selected }">
        <c-action-btn
          type="delete"
          @click="$emit('remove-selected', selected)"
        />
      </template>
      <template #body.prepend="">
        <tr>
          <td v-if="selectable" />
          <td v-if="expandable" />
          <td v-for="header in sortedHeaders" :key="header.value">
            <v-layout class="py-1 gap-2" column>
              <external-data-table-column-data-type-field
                v-field="columns[header.value]"
                :table-separator="separator"
              />
              <v-flex>
                <external-data-table-column-tag-field
                  v-if="header.value !== '_id' && header.value !== 'actions'"
                  v-field="columns[header.value].tag"
                  :disabled="disabled"
                />
              </v-flex>
            </v-layout>
            <span
              v-if="resizingMode"
              :key="`${header.value}.resize`"
              class="table__resize-handler"
              @mousedown.stop.prevent="startColumnResize(header.value)"
              @click.stop=""
            />
          </td>
        </tr>
      </template>
      <template
        v-for="header in sortedHeaders"
        #[`header.${header.value}`]=""
      >
        <span
          :key="`header.${header.value}`"
          :title="header.text"
          :class="{ 'table-cell__header--error': header.errors?.length }"
          class="table-cell__content"
        >
          {{ header.text }}
          <c-help-icon
            v-if="header.errors?.length"
            :text="header.errors.join(', ')"
            icon="help"
          />
        </span>
        <span
          v-if="draggingMode"
          :key="`header.${header.value}.drag`"
          class="table__dragging-handler"
          @click.stop=""
        />
        <span
          v-if="resizingMode"
          :key="`header.cell.${header.value}.resize`"
          class="table__resize-handler"
          @mousedown.stop.prevent="startColumnResize(header.value)"
          @click.stop=""
        />
      </template>
      <template
        v-for="header in sortedHeaders"
        #[header.value]="{ item }"
      >
        <span
          :key="`${header.value}`"
          :title="item[header.value]"
          :class="{ 'table-cell__content--error': item[header.value]?.transform_error }"
          class="table-cell__content"
        >
          {{ item[header.value]?.initial_value ?? item[header.value] }}
        </span>
        <span
          v-if="resizingMode"
          :key="`${header.value}.resize`"
          class="table__resize-handler"
          @mousedown.stop.prevent="startColumnResize(header.value)"
          @click.stop=""
        />
      </template>
      <template #actions="{ item }">
        <v-layout>
          <c-action-btn
            type="edit"
            @click="$emit('edit', item)"
          />
          <c-action-btn
            type="duplicate"
            @click="$emit('duplicate', item)"
          />
          <c-action-btn
            type="delete"
            @click="$emit('remove', item._id)"
          />
        </v-layout>
      </template>
      <template #expand="{ item }">
        <external-data-table-records-list-expand-panel :record="item" />
      </template>
    </c-advanced-data-table>
    <c-alert
      v-if="errorsMessages.length"
      type="error"
    >
      {{ $tc('externalData.fieldsHasError', { count: errorsMessages.length }, errorsMessages.length) }}
      <ul>
        <li v-for="error in errorsMessages" :key="error.name">
          <strong>{{ error.name }}</strong> {{ error.message }}
        </li>
      </ul>
    </c-alert>
  </div>
</template>

<script>
import { isEmpty } from 'lodash';
import { computed, ref, toRef } from 'vue';

import { DENSE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useHTMLElement } from '@/hooks/html-elements';
import { useTableColumnsSettings } from '@/hooks/table/columns-settings';

import ExternalDataTableColumnTagField from '../form/fields/external-data-table-column-tag-field.vue';
import ExternalDataTableColumnDataTypeField from '../form/fields/external-data-table-column-data-type-field.vue';

import ExternalDataTableRecordsListExpandPanel from './external-data-table-records-list-expand-panel.vue';

const DRAGGABLE_CLASS = 'external-data-table-records__draggable-column';

export default {
  components: {
    ExternalDataTableColumnTagField,
    ExternalDataTableColumnDataTypeField,
    ExternalDataTableRecordsListExpandPanel,
  },
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    records: {
      type: Array,
      default: () => [],
    },
    columns: {
      type: Object,
      default: () => ({}),
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    hasStructure: {
      type: Boolean,
      default: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    withActions: {
      type: Boolean,
      default: false,
    },
    withToolbar: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    densable: {
      type: Boolean,
      default: false,
    },
    exportable: {
      type: Boolean,
      default: false,
    },
    draggable: {
      type: Boolean,
      default: false,
    },
    resizable: {
      type: Boolean,
      default: false,
    },
    editableColumns: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Number,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    columnsSettings: {
      type: Object,
      default: () => ({}),
    },
    cellsContentBehavior: {
      type: String,
      required: false,
    },
    separator: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    /**
     * HTML ELEMENTS
     */
    const tableElement = ref(null);

    const { element: tableHeaderElement } = useHTMLElement({
      parentElement: tableElement,
      selector: 'table > thead tr',
    });

    /**
     * COMPUTED
     */
    const isEmptyColumns = computed(() => isEmpty(props.columns));

    const errorsMessages = computed(() => (
      Object.values(props.columns).reduce((acc, column) => {
        if (column.messages?.length) {
          acc.push({
            name: column.name,
            message: `${column.rows.slice(0, 5)?.join(', ')}${column.rows.length > 5 ? t('externalData.andMore') : ''}`,
          });
        }

        return acc;
      }, [])
    ));

    const headers = computed(() => {
      const result = Object.keys(props.columns).map(column => ({
        value: column,
        text: column,
        class: DRAGGABLE_CLASS,
        sortable: true,
        errors: column.messages,
      }));

      if (props.withActions) {
        result.push({
          value: 'actions',
          class: DRAGGABLE_CLASS,
          text: t('common.actionsLabel'),
          sortable: false,
        });
      }

      return result;
    });

    const isSmallDense = computed(() => props.dense === DENSE_TYPES.small);
    const isMediumDense = computed(() => props.dense === DENSE_TYPES.medium);

    const updatePage = page => emit('update:options', { ...props.options, page });

    const {
      tableStyle,
      tableClass,
      isColumnsEditing,
      draggingMode,
      resizingMode,
      sortedHeaders,
      sortedHeadersWithWidth,

      startColumnResize,
      toggleColumnEditingMode,
      resetColumnsSettings,
    } = useTableColumnsSettings({
      tableHeaderElement,
      headers,
      draggable: toRef(props, 'draggable'),
      resizable: toRef(props, 'draggable'),
      columnsSettings: toRef(props, 'columnsSettings'),
      cellsContentBehavior: toRef(props, 'cellsContentBehavior'),
      draggableClass: DRAGGABLE_CLASS,
    });

    return {
      tableElement,

      isEmptyColumns,
      errorsMessages,
      headers,
      isSmallDense,
      isMediumDense,

      updatePage,

      tableStyle,
      tableClass,
      isColumnsEditing,
      draggingMode,
      resizingMode,
      sortedHeaders,
      sortedHeadersWithWidth,

      startColumnResize,
      toggleColumnEditingMode,
      resetColumnsSettings,
    };
  },
};
</script>

<style lang="scss">
.external-data-table {
  &--fixed {
    & > .v-data-table__wrapper > table {
      table-layout: fixed;
      width: var(--external-data-table-width);
      max-width: unset;
    }

    th[data-value="data-table-select"], th[data-value="data-table-expand"] {
      width: 60px !important;
    }
  }

  .table-cell__header--error, .table-cell__content--error  {
    color: var(--v-error-base);
  }

  .table-cell__content--error {
    font-style: italic;
  }

  &-records__draggable-column {
    position: relative;
  }
}
</style>
