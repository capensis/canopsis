<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :disabled="onlyColor"
    />
    <v-text-field
      v-field="form.description"
      v-validate="'required'"
      :label="$t('modals.createPbehaviorType.fields.description')"
      :error-messages="errors.collect('description')"
      :disabled="onlyColor"
      name="description"
    />
    <v-layout justify-space-between>
      <v-flex xs6>
        <pbehavior-type-default-type-field
          :value="form.type"
          :disabled="onlyColor"
          class="mr-2"
          @input="updateDefaultType"
          @update:color="updateColor"
        />
      </v-flex>
      <v-flex
        class="ml-2"
        xs6
      >
        <c-priority-field
          v-field="form.priority"
          :disabled="onlyColor"
          :loading="pendingPriority"
          required
        />
      </v-flex>
    </v-layout>
    <c-icon-field
      v-field="form.icon_name"
      :label="$t('modals.createPbehaviorType.fields.iconName')"
      :hint="$t('modals.createPbehaviorType.iconNameHint')"
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
    <v-flex
      class="mt-2"
      xs12
    >
      <v-alert
        :value="onlyColor"
        color="info"
      >
        {{ $t('pbehavior.types.defaultType') }}
      </v-alert>
    </v-flex>
    <c-enabled-field
      v-field="form.hidden"
      :label="$t('pbehavior.types.hidden')"
    />
    <c-color-picker-field
      v-field="form.color"
      class="mt-2"
      required
      @input="setColorWasChanged"
    />
  </v-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import PbehaviorTypeDefaultTypeField from './pbehavior-type-default-type-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorTypeDefaultTypeField },
  mixins: [formMixin],
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
  data() {
    return {
      colorWasChanged: false,
    };
  },
  methods: {
    setColorWasChanged() {
      this.colorWasChanged = true;
    },

    updateDefaultType(type, color) {
      const newForm = {
        ...this.form,

        type,
      };

      if (!this.colorWasChanged) {
        newForm.color = color;
      }

      this.updateModel(newForm);
    },

    updateColor(color) {
      if (this.form.color) {
        if (this.colorWasChanged) {
          return;
        }

        if (this.form.color !== color) {
          this.colorWasChanged = true;
          return;
        }
      }

      this.updateField('color', color);
    },
  },
};
</script>
