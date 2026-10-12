<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <v-divider />
    <field-periodic-refresh v-model="form.parameters.periodic_refresh" />
    <v-divider />
    <field-storage
      v-model="form.parameters.directory"
      :title="t('settings.resultDirectory')"
      :disabled="form.parameters.is_api"
      name="parameters.directory"
      required
      @add="editResultDirectory"
      @edit="editResultDirectory"
      @remove="removeResultDirectory"
    />
    <v-divider />
    <widget-settings-group :title="t('settings.advancedSettings')">
      <field-switcher
        v-model="form.parameters.is_api"
        :title="t('settings.receiveByApi')"
      />
      <v-divider />
      <field-storages
        v-model="form.parameters.screenshot_directories"
        :disabled="form.parameters.is_api"
        :help-text="t('settings.screenshotDirectories.helpText')"
        :title="t('settings.screenshotDirectories.title')"
        @add="editScreenshotStorage"
        @edit="editScreenshotStorage"
      />
      <v-divider />
      <field-storages
        v-model="form.parameters.video_directories"
        :disabled="form.parameters.is_api"
        :help-text="t('settings.videoDirectories.helpText')"
        :title="t('settings.videoDirectories.title')"
        @add="editVideoStorage"
        @edit="editVideoStorage"
      />
      <v-divider />
      <field-file-name-masks v-model="form.parameters" />
    </widget-settings-group>
    <v-divider />
    <widget-colors-form v-model="form.parameters.colors" />
  </widget-settings>
</template>

<script setup>
import { toRef } from 'vue';

import { MODALS, SIDE_BARS } from '@/constants';

import { uid } from '@/helpers/uid';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useTestSuite } from '@/hooks/store/modules/test-suite';
import { useWidgetSettings } from '@/hooks/widget/settings';

import FieldPeriodicRefresh from '../form/fields/periodic-refresh.vue';
import FieldSwitcher from '../form/fields/switcher.vue';
import FieldTitle from '../form/fields/title.vue';
import WidgetColorsForm from '../form/fields/widget-colors-form.vue';
import WidgetSettingsGroup from '../partials/widget-settings-group.vue';
import WidgetSettings from '../partials/widget-settings.vue';

import FieldFileNameMasks from './form/fields/file-name-masks.vue';
import FieldStorage from './form/fields/storage.vue';
import FieldStorages from './form/fields/storages.vue';

defineOptions({
  name: SIDE_BARS.testingWeatherSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const modals = useModals();
const { validateTestSuiteDirectory } = useTestSuite();
const sidebar = toRef(props, 'sidebar');

const { form, hasChanges, submitting, close, submit } = useWidgetSettings({
  sidebar,
});

const getStorageByString = directory => ({
  directory: directory || '',
  key: uid(),
});

const showDefineStorageModal = ({ title, field, action }) => {
  modals.show({
    name: MODALS.textFieldEditor,
    config: {
      title,
      field: {
        ...field,
        name: 'directory',
        validationRules: { required: true },
      },
      action: async (directory) => {
        await validateTestSuiteDirectory({ data: { directory } });

        await action(directory);
      },
    },
  });
};

const editResultDirectory = (value) => {
  showDefineStorageModal({
    title: t('modals.defineStorage.title'),
    field: {
      value,
      placeholder: t('modals.defineStorage.field.placeholder'),
    },
    action: (directory) => {
      form.value.parameters.directory = directory;
    },
  });
};

const removeResultDirectory = () => {
  form.value.parameters.directory = '';
};

const editScreenshotStorage = (storage, index) => {
  showDefineStorageModal({
    title: t('modals.defineScreenshotStorage.title'),
    field: {
      value: storage && storage.directory,
      placeholder: t('modals.defineScreenshotStorage.field.placeholder'),
    },
    action: (directory) => {
      if (storage) {
        form.value.parameters.screenshot_directories.splice(
          index,
          1,
          getStorageByString(directory),
        );
      } else {
        form.value.parameters.screenshot_directories.push(getStorageByString(directory));
      }
    },
  });
};

const editVideoStorage = (storage, index) => {
  showDefineStorageModal({
    title: t('modals.defineVideoStorage.title'),
    field: {
      value: storage && storage.directory,
      placeholder: t('modals.defineVideoStorage.field.placeholder'),
    },
    action: (directory) => {
      if (storage) {
        form.value.parameters.video_directories.splice(index, 1, getStorageByString(directory));
      } else {
        form.value.parameters.video_directories.push(getStorageByString(directory));
      }
    },
  });
};

defineExpose({
  close,
});
</script>
