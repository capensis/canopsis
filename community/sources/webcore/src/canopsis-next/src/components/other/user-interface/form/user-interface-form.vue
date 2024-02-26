<template>
  <v-layout column>
    <v-layout>
      <v-text-field
        v-field="form.app_title"
        :disabled="disabled"
        :label="$t('userInterface.appTitle')"
      />
    </v-layout>
    <c-duration-field
      v-field="form.popup_timeout.info"
      :label="$t('userInterface.infoPopupTimeout')"
      name="popup_timeout.info"
    />
    <c-duration-field
      v-field="form.popup_timeout.error"
      :label="$t('userInterface.errorPopupTimeout')"
      name="popup_timeout.error"
    />
    <v-layout>
      <c-language-field
        v-field="form.language"
        :label="$t('userInterface.language')"
      />
    </v-layout>
    <v-layout>
      <c-number-field
        v-field="form.max_matched_items"
        :label="$t('userInterface.maxMatchedItems')"
        :min="1"
        name="max_matched_items"
      >
        <template #append="">
          <c-help-icon
            :text="$t('userInterface.tooltips.maxMatchedItems')"
            color="grey darken-1"
            icon="help"
            left
          />
        </template>
      </c-number-field>
    </v-layout>
    <v-layout>
      <c-number-field
        v-field="form.check_count_request_timeout"
        :label="$t('userInterface.checkCountRequestTimeout')"
        :min="1"
        name="check_count_request_timeout"
      >
        <template #append="">
          <c-help-icon
            :text="$t('userInterface.tooltips.checkCountRequestTimeout')"
            color="grey darken-1"
            icon="help"
            left
          />
        </template>
      </c-number-field>
    </v-layout>
    <v-layout>
      <c-timezone-field
        v-field="form.timezone"
        disabled
      />
    </v-layout>
    <v-layout>
      <v-flex xs6>
        <c-enabled-field
          v-field="form.allow_change_severity_to_info"
          :label="$t('userInterface.allowChangeSeverityToInfo')"
        />
      </v-flex>
      <v-flex xs6>
        <c-enabled-field
          v-field="form.show_header_on_kiosk_mode"
          :label="$t('userInterface.showHeaderOnKioskMode')"
        />
      </v-flex>
    </v-layout>
    <v-layout>
      <v-flex>
        <text-editor-field
          v-field="form.footer"
          :label="$t('userInterface.footer')"
          :config="textEditorConfig"
          public
        />
      </v-flex>
    </v-layout>
    <v-layout class="mt-3">
      <v-flex>
        <text-editor-field
          v-field="form.login_page_description"
          :label="$t('userInterface.description')"
          :config="textEditorConfig"
          public
        />
      </v-flex>
    </v-layout>
    <v-layout class="mt-3">
      <v-flex>
        <span class="v-label file-selector__label">{{ $t('userInterface.logo') }}</span>
        <v-layout>
          <file-selector
            ref="fileSelector"
            :max-file-size="maxFileSize"
            :disabled="disabled"
            class="mt-1"
            accept="image/*"
            name="logo"
            with-files-list
            @change="changeLogoFile"
          />
        </v-layout>
      </v-flex>
    </v-layout>
  </v-layout>
</template>

<script>
import { MAX_ICON_SIZE_IN_KB } from '@/constants';

import { formMixin } from '@/mixins/form';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: {
    FileSelector,
    TextEditorField,
  },
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    maxFileSize() {
      return MAX_ICON_SIZE_IN_KB;
    },

    textEditorConfig() {
      return { disabled: this.disabled };
    },
  },
  methods: {
    async changeLogoFile([file] = []) {
      this.updateField('logo', file);
    },

    reset() {
      this.$refs.fileSelector.clear();
    },
  },
};
</script>

<style lang="scss" scoped>
.file-selector {
  &__label {
    font-size: .85em;
    display: block;
  }
}
</style>
