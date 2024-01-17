<template>
  <file-selector
    class="my-2"
    v-validate="'required'"
    :error-messages="errors.collect('file')"
    name="file"
    accept=".svg"
    hide-details
    @change="chooseIcon"
  >
    <template #activator="{ clear, on, ...attrs }">
      <v-layout
        class="cursor-default"
        v-if="file"
        align-center
      >
        <v-flex xs1>
          <v-icon
            large
          >
            $vuetify.icons.{{ uniqueIconName }}
          </v-icon>
        </v-flex>
        <v-flex xs10>
          {{ file.name }}
        </v-flex>
        <v-flex xs1>
          <c-action-btn
            type="delete"
            @click.prevent="clear"
          />
        </v-flex>
      </v-layout>
      <v-btn
        class="import-btn ma-0"
        v-else
        v-bind="attrs"
        v-on="on"
        color="primary"
      >
        <v-icon left>
          file_upload
        </v-icon>
        <span>{{ $t('common.chooseFile') }}</span>
      </v-btn>
    </template>
  </file-selector>
</template>

<script>
import { omit } from 'lodash';

import { uid } from '@/helpers/uid';
import { normalizeHtml } from '@/helpers/html';
import { getFileTextContent } from '@/helpers/file/file-select';

import { formBaseMixin } from '@/mixins/form';

import FileSelector from '@/components/forms/fields/file-selector.vue';

export default {
  inject: ['$validator'],
  components: { FileSelector },
  mixins: [formBaseMixin],
  model: {
    prop: 'file',
    event: 'change',
  },
  props: {
    file: {
      type: File,
      default: null,
    },
  },
  data() {
    return {
      uniqueIconName: uid(),
    };
  },
  created() {
    this.parseIcon(this.file);
  },
  beforeDestroy() {
    this.removeIconFromVuetify();
  },
  methods: {
    async parseIcon(file) {
      if (!file) {
        this.removeIconFromVuetify();

        return;
      }

      const template = await getFileTextContent(file);

      this.addIconIntoVuetify(template);
    },

    addIconIntoVuetify(template) {
      this.$vuetify.icons.values = {
        ...this.$vuetify.icons.values,

        [this.uniqueIconName]: {
          component: {
            name: this.uniqueIconName,
            template: normalizeHtml(template),
          },
        },
      };
    },

    removeIconFromVuetify() {
      if (!this.$vuetify.icons.values[this.uniqueIconName]) {
        return;
      }

      this.$vuetify.icons.values = omit(this.$vuetify.icons.values, [this.uniqueIconName]);
    },

    chooseIcon([file = null] = []) {
      this.parseIcon(file);
      this.updateModel(file);
    },
  },
};
</script>
