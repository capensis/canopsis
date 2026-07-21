<template>
  <v-layout
    class="gap-3"
    column
  >
    <c-enabled-field
      v-field="form.visible"
      :label="$t('pbehavior.visible')"
      with-background
    />
    <c-name-field
      v-field="form.name"
      :disabled="onlyColor"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('modals.createPbehaviorType.fields.description')">
        <c-description-field
          v-field="form.description"
          :disabled="onlyColor"
          name="description"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createPbehaviorType.fields.type')">
        <pbehavior-type-default-type-field
          :value="form.type"
          :disabled="onlyColor"
          @input="updateDefaultType"
          @update:color="updateColor"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createPbehaviorType.fields.priority')">
        <c-priority-field
          v-field="form.priority"
          :disabled="onlyColor"
          :loading="pendingPriority"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createPbehaviorType.fields.iconName')">
        <c-icon-field
          v-field="form.icon_name"
          :placeholder="$t('modals.createPbehaviorType.iconNamePlaceholder')"
          :disabled="onlyColor"
          :required="!onlyColor"
        >
          <template #no-data="">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title v-html="$t('modals.createPbehaviorType.errors.iconName')" />
              </v-list-item-content>
            </v-list-item>
          </template>
        </c-icon-field>
      </c-form-block-row>

      <c-form-block-row
        :label="$t('common.color')"
        align-center
      >
        <v-layout class="gap-2" align-center justify-start>
          <span v-if="onlyColor" class="text--secondary">{{ $t('pbehavior.types.defaultType') }}</span>
          <c-color-picker-field
            v-field="form.color"
            :justify-end="onlyColor"
            required
            @input="setColorWasChanged"
          />
        </v-layout>
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import PbehaviorTypeDefaultTypeField from './pbehavior-type-default-type-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorTypeDefaultTypeField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    onlyColor: {
      type: Boolean,
      default: false,
    },
    pendingPriority: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);

    const colorWasChanged = ref(false);

    /**
     * Marks the color as manually edited so it is no longer overwritten when the default type changes.
     */
    const setColorWasChanged = () => colorWasChanged.value = true;

    /**
     * Updates the selected canonical type and syncs the form color from the default type when the
     * user has not changed the color manually.
     *
     * @param {string} type - Canonical pbehavior type value.
     * @param {string} color - Default color associated with the selected type.
     */
    const updateDefaultType = (type, color) => {
      const newForm = {
        ...props.form,

        type,
      };

      if (!colorWasChanged.value) {
        newForm.color = color;
      }

      updateModel(newForm);
    };

    /**
     * Applies the color suggested by the default type field. Ignores updates once the user has
     * manually changed the color, and marks the color as user-edited when it differs from the form value.
     *
     * @param {string} color - Default color emitted by the selected canonical type.
     */
    const updateColor = (color) => {
      if (props.form.color) {
        if (colorWasChanged.value) {
          return;
        }

        if (props.form.color !== color) {
          colorWasChanged.value = true;

          return;
        }
      }

      updateField('color', color);
    };

    return {
      setColorWasChanged,
      updateDefaultType,
      updateColor,
    };
  },
};
</script>
