<template>
  <file-selector
    ref="fileSelector"
    multiple
    hide-details
    @change="changeFiles"
  >
    <template #activator="{ on: fileSelectorOn }">
      <v-tooltip top>
        <template #activator="{ on: tooltipOn }">
          <v-btn
            :loading="pending"
            color="indigo"
            small
            dark
            fab
            v-on="{ ...fileSelectorOn, ...tooltipOn }"
          >
            <v-icon small>
              cloud_upload
            </v-icon>
          </v-btn>
        </template>
        <span>{{ $t('snmpRule.uploadMib') }}</span>
      </v-tooltip>
    </template>
  </file-selector>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { getFileTextContent } from '@/helpers/file/file-select';

import FileSelector from '@/components/forms/fields/file-selector.vue';

const { mapActions } = createNamespacedHelpers('snmpMib');

export default {
  components: { FileSelector },
  data() {
    return {
      pending: false,
    };
  },
  methods: {
    ...mapActions({
      uploadSnmpMib: 'upload',
    }),

    async changeFiles(files = []) {
      try {
        this.pending = true;

        if (files.length) {
          const promises = files.sort((a, b) => {
            const matchesA = a.name.match(/^(\d+)/) || ['0'];
            const matchesB = b.name.match(/^(\d+)/) || ['0'];

            return parseInt(matchesA[0], 10) - parseInt(matchesB[0], 10);
          }).map(file => getFileTextContent(file));

          const results = await Promise.all(promises);

          const fileContent = results.map((content, index) => ({
            filename: files[index].name,
            data: content,
          }));

          const { counts } = await this.uploadSnmpMib({
            data: { filecontent: fileContent },
          });

          this.$popups.success({
            text: this.$tc('snmpRule.uploadedMibPopup', files.length, counts),
            autoClose: 10000,
          });
        }
      } catch (err) {
        console.warn(err);

        this.$popups.error({
          text: err.error ?? this.$t('errors.default'),
          autoClose: 10000,
        });
      } finally {
        this.pending = false;
        this.$refs.fileSelector.clear();
      }
    },
  },
};
</script>
