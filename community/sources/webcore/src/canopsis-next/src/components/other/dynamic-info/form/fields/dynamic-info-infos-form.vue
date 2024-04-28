<template>
  <div>
    <v-alert
      :value="errors.has(name)"
      type="error"
    >
      {{ errorMessages }}
    </v-alert>
    <v-layout justify-end>
      <v-btn
        class="primary"
        fab
        small
        text
        @click="showAddInfoModal"
      >
        <v-icon>add</v-icon>
      </v-btn>
    </v-layout>
    <v-data-table
      :headers="headers"
      :items="form"
      :no-data-text="$t('common.noData')"
      item-key="key"
    >
      <template #item="{ item }">
        <tr>
          <td>{{ item.name }}</td>
          <td>{{ item.value }}</td>
          <td>
            <v-layout>
              <c-action-btn
                type="edit"
                @click="showEditInfoModal(item)"
              />
              <c-action-btn
                type="delete"
                @click="removeInfo(item)"
              />
            </v-layout>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { formArrayMixin, formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin, formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'infos',
    },
  },
  computed: {
    errorMessages() {
      return this.errors.collect(this.name, undefined, false)
        ?.map(({ rule, msg }) => {
          const customMessage = {
            required: this.$t('modals.createDynamicInfo.errors.emptyInfos'),
          }[rule];

          return customMessage || msg;
        })
        .join('\n');
    },

    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.value'), value: 'value' },
        { text: this.$t('common.actionsLabel'), value: 'actions' },
      ];
    },
  },
  watch: {
    form() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.attachRequiredRule();
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    attachRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => !!this.form.length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
    },

    findInfoIndex(info = {}) {
      return this.form.findIndex(({ name }) => name === info.name);
    },

    showAddInfoModal() {
      this.$modals.show({
        name: MODALS.createDynamicInfoInformation,
        config: {
          existingNames: this.form.map(info => info.name),
          action: newInfo => this.addItemIntoArray(newInfo),
        },
      });
    },

    showEditInfoModal(info) {
      this.$modals.show({
        name: MODALS.createDynamicInfoInformation,
        config: {
          info,

          existingNames: this.form.map(({ name }) => name),
          action: (newInfo) => {
            const index = this.findInfoIndex(info);

            this.updateItemInArray(index, newInfo);
          },
        },
      });
    },

    removeInfo(info) {
      const index = this.findInfoIndex(info);

      this.removeItemFromArray(index);
    },
  },
};
</script>
