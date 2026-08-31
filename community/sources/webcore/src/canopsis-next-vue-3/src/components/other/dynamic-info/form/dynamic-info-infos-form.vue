<template>
  <c-layout class="py-3" gap="3" column>
    <v-card v-for="(field, index) in infosField" :key="field.key">
      <v-card-text>
        <c-layout gap="2" column>
          <c-layout align="center" gap="2">
            <c-name-field
              v-model="field.name"
              :label="$t('common.name')"
              :name="`infos[${index}].name`"
              required
            />
            <c-action-btn v-if="infosField.length > 1" type="delete" @click="remove(index)" />
          </c-layout>

          <c-select
            v-model="field.type"
            :items="typeItems"
            :label="$t('common.type')"
            :name="`infos[${index}].type`"
            required
          />

          <c-text-field
            v-model="field.value"
            :label="$t('common.value')"
            :name="`infos[${index}].value`"
            required
          />
        </c-layout>
      </v-card-text>
    </v-card>

    <c-layout gap="2">
      <v-btn color="primary" variant="outlined" @click="addItem">
        {{ $t('modals.createDynamicInfo.infosSection.addInfos') }}
      </v-btn>

      <v-btn color="primary" variant="outlined" @click="showAddInfosFromTemplateModal">
        {{ $t('modals.createDynamicInfo.infosSection.addInfosFromTemplate') }}
      </v-btn>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { DYNAMIC_INFO_INFORMATION_TYPES, MODALS } from '@/constants';

import { dynamicInfoInformationToForm } from '@/helpers/entities/dynamic-info/information/form';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const form = useFormContext();
    const infosField = form.array('infos');

    const typeItems = computed(() => Object.values(DYNAMIC_INFO_INFORMATION_TYPES).map(type => ({
        title: t(`dynamicInfo.informationTypes.${type}`),
        value: type,
      })),
    );

    const addItem = () => infosField.value.push(dynamicInfoInformationToForm());
    const remove = index => infosField.value.splice(index, 1);
    const showAddInfosFromTemplateModal = () => modals.show({
        name: MODALS.addDynamicInfoInfosFromTemplate,
        config: {
          action: template => template.names.forEach(name => infosField.value.push(dynamicInfoInformationToForm({ name })),
            ),
        },
      });

    return {
      infosField,
      typeItems,
      addItem,
      showAddInfosFromTemplateModal,
      remove,
    };
  },
};
</script>
