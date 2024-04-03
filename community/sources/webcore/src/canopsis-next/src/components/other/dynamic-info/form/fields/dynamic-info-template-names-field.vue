<template>
  <div>
    <slot
      v-if="!names.length"
      name="no-data"
    />
    <v-layout
      v-for="(name, index) in names"
      :key="name.key"
      justify-space-between
    >
      <v-flex xs11>
        <v-text-field
          v-field="names[index].value"
          v-validate="'required'"
          :error-messages="errors.collect(`name[${name.key}]`)"
          :name="`name[${name.key}]`"
          :placeholder="$t('common.name')"
        />
      </v-flex>
      <v-flex xs1>
        <v-btn
          color="error"
          icon
          @click="removeItemFromArray(index)"
        >
          <v-icon>delete</v-icon>
        </v-btn>
      </v-flex>
    </v-layout>
    <v-btn
      class="primary mx-0"
      @click="showAddValueModal"
    >
      {{ $t('modals.createDynamicInfoTemplate.buttons.addName') }}
    </v-btn>
    <v-alert
      :value="errors.has('names')"
      type="error"
    >
      <span>{{ $t('modals.createDynamicInfoTemplate.errors.noNames') }}</span>
    </v-alert>
  </div>
</template>

<script>
import { generateTemplateFormName } from '@/helpers/entities/dynamic-info/template/form';

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'names',
    event: 'input',
  },
  props: {
    names: {
      type: Array,
      default: () => [],
    },
  },
  created() {
    this.$validator.attach({
      name: 'names',
      rules: 'required:true',
      getter: () => this.names.length > 0,
      vm: this,
    });
  },
  methods: {
    showAddValueModal() {
      this.addItemIntoArray(generateTemplateFormName());

      this.$nextTick(() => this.$validator.validate('names'));
    },
  },
};
</script>
