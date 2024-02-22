<template>
  <v-select
    v-validate
    :value="value"
    :label="label || $t('role.selectTemplate')"
    :loading="pending"
    :items="preparedItems"
    :error-messages="errors.collect(name)"
    :name="name"
    :disabled="disabled"
    class="role-template-field"
    item-text="name"
    item-value="permissions"
    return-object
    clearable
    @input="updatePermissions"
  >
    <template #item="{ item }">
      <v-layout justify-space-between>
        {{ item.name }}<span class="role-template-field__item-description">{{ item.description }}</span>
      </v-layout>
    </template>
  </v-select>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

import { rolePermissionsToForm } from '@/helpers/entities/role/form';

import { formBaseMixin } from '@/mixins/form/base';

const { mapActions: mapRoleActions } = createNamespacedHelpers('role');

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'template',
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      items: [],
      pending: false,
    };
  },
  computed: {
    preparedItems() {
      return this.items.map(template => ({
        ...template,
        permissions: rolePermissionsToForm(template.permissions),
      }));
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapRoleActions({ fetchRoleTemplatesListWithoutStore: 'fetchTemplatesListWithoutStore' }),

    updatePermissions(value) {
      this.updateModel(value ? value.permissions : {});
    },

    async fetchList() {
      this.pending = true;

      const { data: roleTemplates } = await this.fetchRoleTemplatesListWithoutStore({
        params: {
          types: this.types,
          limit: MAX_LIMIT,
        },
      });

      this.items = roleTemplates;
      this.pending = false;
    },
  },
};
</script>

<style lang="scss">
.role-template-field {
  &__item-description {
    opacity: 0.7;
  }
}
</style>
