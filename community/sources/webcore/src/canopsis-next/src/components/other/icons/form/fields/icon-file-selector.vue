<template>
  <file-drag-selector
    :file-type-label="$t('common.fileSelector.fileTypes.svg')"
    accept=".svg,.svg+xml"
    required
    @change="chooseIcon"
  >
    <template #selected="{ clear }">
      <v-layout
        class="cursor-default"
        align-center
      >
        <v-flex xs2>
          <v-icon large>
            $vuetify.icons.{{ uniqueIconName }}
          </v-icon>
        </v-flex>
        <v-flex xs9>
          {{ file.name }}
        </v-flex>
        <v-flex
          class="text-right"
          xs2
        >
          <c-action-btn
            type="delete"
            @click.prevent="clear"
          />
        </v-flex>
      </v-layout>
    </template>
  </file-drag-selector>
</template>

<script>
import { uid } from '@/helpers/uid';
import { getFileTextContent } from '@/helpers/file/file-select';

import { formBaseMixin } from '@/mixins/form';
import { vuetifyCustomIconsBaseMixin } from '@/mixins/vuetify/custom-icons/base';

import FileDragSelector from '@/components/forms/fields/file-drag-selector.vue';

export default {
  components: { FileDragSelector },
  mixins: [
    formBaseMixin,
    vuetifyCustomIconsBaseMixin,
  ],
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
      uniqueIconName: uid('icon'),
    };
  },
  created() {
    this.parseIcon(this.file);
  },
  beforeDestroy() {
    this.unregisterIconFromVuetify(this.uniqueIconName);
  },
  methods: {
    async parseIcon(file) {
      if (!file) {
        this.unregisterIconFromVuetify(this.uniqueIconName);

        return;
      }

      const template = await getFileTextContent(file);

      this.registerIconInVuetify(this.uniqueIconName, template);
    },

    chooseIcon([file = null] = []) {
      this.parseIcon(file);
      this.updateModel(file);
    },
  },
};
</script>
