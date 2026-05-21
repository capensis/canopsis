<template>
  <v-layout
    class="gap-3"
    column
  >
    <c-name-field
      v-field="form.name"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('pbehavior.pbehaviorType')">
        <c-pbehavior-type-field
          v-field="form.type"
          name="type"
          required
        />
      </c-form-block-row>

      <c-form-block-row
        :label="$t('common.file')"
        align-center
      >
        <file-selector
          name="file"
          required
          with-files-list
          @change="changeFiles"
        >
          <template #activator="{ on, disabled }">
            <v-tooltip top>
              <template #activator="{ on: tooltipOn }">
                <v-btn
                  :color="errors.has('file') ? 'error' : 'primary'"
                  :disabled="disabled"
                  small
                  outlined
                  v-on="{ ...on, ...tooltipOn }"
                >
                  <v-icon>cloud_upload</v-icon>
                </v-btn>
              </template>
              <span>{{ $t('common.chooseFile') }}</span>
            </v-tooltip>
          </template>
        </file-selector>
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { onMounted } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import FileSelector from '@/components/forms/fields/file-selector.vue';

export default {
  inject: ['$validator'],
  components: { FileSelector },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    /**
     * Stores the first selected file in the import form.
     *
     * @param {File[]} [files=[]] - Selected files from the file selector.
     */
    const changeFiles = (files = []) => updateField('file', files[0]);

    onMounted(fetchPbehaviorTypesFieldList);

    return {
      changeFiles,
    };
  },
};
</script>
