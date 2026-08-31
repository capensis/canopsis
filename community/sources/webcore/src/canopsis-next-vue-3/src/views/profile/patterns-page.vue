<template>
  <div>
    <c-page-header>{{ $t('pattern.patterns') }}</c-page-header>

    <v-card class="ma-2">
      <template v-if="hasReadAnyCorporatePatternAccess">
        <v-tabs v-model="activeTab" color="primary" fixed-tabs>
          <v-tab :value="PATTERN_TABS.patterns">
            {{ $t('pattern.myPatterns') }}
          </v-tab>
          <v-tab :value="PATTERN_TABS.corporatePatterns">
            {{ $t('pattern.corporatePatterns') }}
          </v-tab>
        </v-tabs>

        <v-window v-model="activeTab">
          <v-window-item :value="PATTERN_TABS.patterns">
            <patterns
              @edit="showEditPatternModal"
              @remove="showDeletePatternModal"
              @remove-selected="showDeleteSelectedPatternsModal"
            />
          </v-window-item>
          <v-window-item :value="PATTERN_TABS.corporatePatterns">
            <corporate-patterns
              @edit="showEditPatternModal"
              @remove="showDeletePatternModal"
              @remove-selected="showDeleteSelectedPatternsModal"
            />
          </v-window-item>
        </v-window>
      </template>

      <patterns
        v-else
        @edit="showEditPatternModal"
        @remove="showDeletePatternModal"
        @remove-selected="showDeleteSelectedPatternsModal"
      />
    </v-card>

    <c-fab-expand-btn :has-access="hasAccessToCreatePattern" @refresh="refresh">
      <c-action-fab-btn
        :tooltip="createPbehaviorTitle"
        icon="pause"
        color="grey-darken-1"
        top
        @click.stop="showCreatePbehaviorPatternModal"
      />
      <c-action-fab-btn
        :tooltip="createEntityTitle"
        icon="perm_identity"
        color="blue"
        top
        @click.stop="showCreateEntityPatternModal"
      />
      <c-action-fab-btn
        :tooltip="createAlarmTitle"
        icon="notification_important"
        color="error"
        top
        @click.stop="showCreateAlarmPatternModal"
      />
      <c-action-fab-btn
        :tooltip="createServiceWeatherTitle"
        icon="view_module"
        color="secondary-lighten-2"
        top
        @click.stop="showCreateServiceWeatherPatternModal"
      />
    </c-fab-expand-btn>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS, PATTERN_TABS, PATTERN_TYPES, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePattern } from '@/hooks/store/modules/pattern';
import { useCorporatePattern } from '@/hooks/store/modules/pattern-corporate';

import CorporatePatterns from '@/components/other/pattern/corporate-patterns.vue';
import Patterns from '@/components/other/pattern/patterns.vue';

export default {
  components: {
    Patterns,
    CorporatePatterns,
  },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const activeTab = ref(PATTERN_TABS.patterns);

    const {
      createPattern,
      updatePattern,
      removePattern,
      bulkRemovePatterns,
      fetchPatternsListWithPreviousParams,
    } = usePattern();
    const { fetchCorporatePatternsListWithPreviousParams } = useCorporatePattern();
    const {
      hasReadAccess: hasReadAnyCorporatePatternAccess,
      hasCreateAccess: hasCreateAnyCorporatePatternAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.profile.corporatePattern);

    const isCorporatePatternsTab = computed(
      () => activeTab.value === PATTERN_TABS.corporatePatterns,
    );
    const hasAccessToCreatePattern = computed(
      () => !isCorporatePatternsTab.value || hasCreateAnyCorporatePatternAccess.value,
    );

    const createAlarmTitle = computed(() => isCorporatePatternsTab.value
        ? t('modals.createCorporateAlarmPattern.create.title')
        : t('modals.createAlarmPattern.create.title'),
    );
    const createEntityTitle = computed(() => isCorporatePatternsTab.value
        ? t('modals.createCorporateEntityPattern.create.title')
        : t('modals.createEntityPattern.create.title'),
    );
    const createPbehaviorTitle = computed(() => isCorporatePatternsTab.value
        ? t('modals.createCorporatePbehaviorPattern.create.title')
        : t('modals.createPbehaviorPattern.create.title'),
    );
    const createServiceWeatherTitle = computed(() => isCorporatePatternsTab.value
        ? t('modals.createCorporateServiceWeatherPattern.create.title')
        : t('modals.createServiceWeatherPattern.create.title'),
    );

    const refresh = () => isCorporatePatternsTab.value
        ? fetchCorporatePatternsListWithPreviousParams()
        : fetchPatternsListWithPreviousParams();

    const showCreatePatternModal = ({ title, type }) => {
      modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: { is_corporate: isCorporatePatternsTab.value },
          title,
          type,
          action: async (pattern) => {
            const result = await createPattern({ data: pattern });

            refresh();

            return result;
          },
        },
      });
    };

    const showCreateAlarmPatternModal = () => showCreatePatternModal({
        title: createAlarmTitle.value,
        type: PATTERN_TYPES.alarm,
      });
    const showCreateEntityPatternModal = () => showCreatePatternModal({
        title: createEntityTitle.value,
        type: PATTERN_TYPES.entity,
      });
    const showCreatePbehaviorPatternModal = () => showCreatePatternModal({
        title: createPbehaviorTitle.value,
        type: PATTERN_TYPES.pbehavior,
      });
    const showCreateServiceWeatherPatternModal = () => showCreatePatternModal({
        title: createServiceWeatherTitle.value,
        type: PATTERN_TYPES.serviceWeather,
      });

    const getEditPatternModalTitle = (pattern) => {
      if (pattern.is_corporate) {
        return {
          [PATTERN_TYPES.alarm]: t('modals.createCorporateAlarmPattern.edit.title'),
          [PATTERN_TYPES.entity]: t('modals.createCorporateEntityPattern.edit.title'),
          [PATTERN_TYPES.pbehavior]: t('modals.createCorporatePbehaviorPattern.edit.title'),
          [PATTERN_TYPES.serviceWeather]: t(
            'modals.createCorporateServiceWeatherPattern.edit.title',
          ),
        }[pattern.type];
      }

      return {
        [PATTERN_TYPES.alarm]: t('modals.createAlarmPattern.edit.title'),
        [PATTERN_TYPES.entity]: t('modals.createEntityPattern.edit.title'),
        [PATTERN_TYPES.pbehavior]: t('modals.createPbehaviorPattern.edit.title'),
        [PATTERN_TYPES.serviceWeather]: t('modals.createServiceWeatherPattern.edit.title'),
      }[pattern.type];
    };

    const showEditPatternModal = (editablePattern) => {
      modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: editablePattern,
          type: editablePattern.type,
          title: getEditPatternModalTitle(editablePattern),
          action: async (pattern) => {
            const result = await updatePattern({ id: editablePattern._id, data: pattern });

            refresh();

            return result;
          },
        },
      });
    };

    const showDeletePatternModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removePattern({ id });

            return refresh();
          },
        },
      });
    };

    const showDeleteSelectedPatternsModal = (selected) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await bulkRemovePatterns({
              data: selected.map(({ _id: id }) => ({ _id: id })),
            });

            return refresh();
          },
        },
      });
    };

    return {
      PATTERN_TABS,
      activeTab,
      hasReadAnyCorporatePatternAccess,
      hasAccessToCreatePattern,
      createAlarmTitle,
      createEntityTitle,
      createPbehaviorTitle,
      createServiceWeatherTitle,
      refresh,
      showCreateAlarmPatternModal,
      showCreateEntityPatternModal,
      showCreatePbehaviorPatternModal,
      showCreateServiceWeatherPatternModal,
      showEditPatternModal,
      showDeletePatternModal,
      showDeleteSelectedPatternsModal,
    };
  },
};
</script>
