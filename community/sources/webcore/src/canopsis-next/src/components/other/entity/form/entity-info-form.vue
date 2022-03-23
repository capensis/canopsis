<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required|unique-name'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-text-field(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('common.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(row)
      v-switch(
        v-model="isListValue",
        :label="$t('context.entityInfo.valueAsList')",
        color="primary"
      )
    template(v-if="isListValue")
      v-layout(v-for="(item, index) in form.value", :key="item.key", row)
        v-flex
          c-mixed-field(
            v-field="form.value[index].value",
            :label="index === 0 ? $t('common.value') : ''",
            required
          )
        v-btn.mx-0(icon, @click="removeListItem(index)")
          v-icon(color="error") delete
      v-layout(row)
        v-btn.ml-0(
          color="primary",
          outline,
          @click="addListItem"
        ) {{ $t('common.add') }}
    v-layout(v-else, row)
      v-flex
        c-mixed-field(
          v-field="form.value",
          :label="$t('common.value')",
          required
        )
</template>

<script>
import { get, isArray } from 'lodash';

import uid from '@/helpers/uid';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
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
    entityInfo: {
      type: Object,
      default: () => ({}),
    },
    infos: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    isListValue: {
      get() {
        return isArray(this.form.value);
      },
      set(value) {
        if (value && !isArray(this.form.value)) {
          this.updateField('value', [{ key: uid(), value: this.form.value }]);
        } else if (!value && isArray(this.form.value)) {
          this.updateField('value', get(this.form.value, [0, 'value'], ''));
        }
      },
    },

    infosNames() {
      return this.infos.map(({ name }) => name);
    },
  },
  created() {
    this.createUniqueValidationRule();
  },
  methods: {
    removeListItem(index) {
      const newValue = this.form.value.filter((item, itemIndex) => index !== itemIndex);

      this.updateField('value', newValue);
    },

    addListItem() {
      const newValue = [...this.form.value, { key: uid(), value: '' }];

      this.updateField('value', newValue);
    },

    createUniqueValidationRule() {
      this.$validator.extend('unique-name', {
        getMessage: () => this.$t('validator.unique'),
        validate: (value) => {
          if (this.entityInfo && this.entityInfo.name === value) {
            return true;
          }

          return !this.infosNames.includes(value);
        },
      });
    },
  },
};
</script>
