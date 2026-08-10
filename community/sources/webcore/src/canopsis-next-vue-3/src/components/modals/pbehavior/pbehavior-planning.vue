<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ $t('modals.pbehaviorPlanning.title') }}</span>
      </template>

      <template #text>
        <pbehavior-planning-calendar
          :pbehaviors-by-id="form.pbehaviorsById"
          :added-pbehaviors-by-id="form.addedPbehaviorsById"
          :changed-pbehaviors-by-id="form.changedPbehaviorsById"
          :removed-pbehaviors-by-id="form.removedPbehaviorsById"
          :with-inherited="config.withInherited"
          :read-only="readOnly"
          :entity-pattern="entityPattern"
          :default-fields="defaultFields"
          :timezone="calendarTimezone"
          @update:pbehaviors-by-id="updateFormField('pbehaviorsById', $event)"
          @update:added-pbehaviors-by-id="updateFormField('addedPbehaviorsById', $event)"
          @update:changed-pbehaviors-by-id="updateFormField('changedPbehaviorsById', $event)"
          @update:removed-pbehaviors-by-id="updateFormField('removedPbehaviorsById', $event)"
          @update:timezone="calendarTimezone = $event"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ $t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { keyBy, omit } from 'lodash';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';
import {
  getPbehaviorNameByEntities,
  getPbehaviorsInitialTimezone,
  pbehaviorToDuplicateForm,
  pbehaviorToRequest,
} from '@/helpers/entities/pbehavior/form';
import { useModals } from '@/hooks/modals';
import { useInfo } from '@/hooks/store/modules/info';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

export default {
  components: { PbehaviorPlanningCalendar },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    const { pbehaviors = [], pbehaviorsToAdd = [] } = this.modal.config ?? {};

    return {
      submitting: false,
      calendarTimezone: getPbehaviorsInitialTimezone(pbehaviors),
      form: {
        pbehaviorsById: keyBy(pbehaviors, '_id'),
        addedPbehaviorsById: keyBy(pbehaviorsToAdd.map(pbehaviorToDuplicateForm), '_id'),
        changedPbehaviorsById: {},
        removedPbehaviorsById: {},
      },
    };
  },

  setup() {
    const modals = useModals();
    const { autoSuggestPbehaviorName } = useInfo();
    const {
      bulkCreatePbehaviors,
      bulkUpdatePbehaviors,
      removePbehaviors,
      updateSeveralPbehaviorComments,
    } = usePbehavior();

    return {
      modals,
      autoSuggestPbehaviorName,
      bulkCreatePbehaviors,
      bulkUpdatePbehaviors,
      removePbehaviors,
      updateSeveralPbehaviorComments,
    };
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },

    defaultName() {
      return this.autoSuggestPbehaviorName
        ? getPbehaviorNameByEntities(this.config.entities, this.calendarTimezone)
        : '';
    },

    defaultFields() {
      const fields = {};

      if (this.defaultName) {
        fields.name = this.defaultName;
      }

      if (this.config.defaultInherited) {
        fields.inherited = this.config.defaultInherited;
      }

      return fields;
    },

    readOnly() {
      return !!this.config.readOnly;
    },

    entityPattern() {
      return this.config.entityPattern;
    },

    isDisabled() {
      return this.submitting || this.readOnly;
    },
  },

  methods: {
    closeModal() {
      this.modals.hide({ id: this.modal.id });
    },

    updateFormField(field, value) {
      this.form = {
        ...this.form,
        [field]: value,
      };
    },

    async createPbehaviorsWithComments(pbehaviors) {
      const response = await this.bulkCreatePbehaviors({ data: pbehaviors });

      await Promise.all(
        response.map(({ id, errors, item: pbehavior }) => {
          if (!errors) {
            return this.updateSeveralPbehaviorComments({
              comments: pbehavior.comments,
              pbehavior: {
                ...pbehavior,
                _id: id,
                comments: [],
              },
            });
          }

          return Promise.reject(errors);
        }),
      );

      return response;
    },

    async updatePbehaviorsWithComments(pbehaviors = [], originalPbehaviors = []) {
      const response = await this.bulkUpdatePbehaviors({ data: pbehaviors });

      await Promise.all(
        pbehaviors.map(pbehavior =>
          this.updateSeveralPbehaviorComments({
            pbehavior: originalPbehaviors.find(({ _id }) => _id === pbehavior._id),
            comments: pbehavior.comments,
          }),
        ),
      );

      return response;
    },

    async submit() {
      this.submitting = true;

      try {
        const { pbehaviors: originalPbehaviors = [] } = this.config;
        const createdPbehaviors = Object.values(this.form.addedPbehaviorsById).map(pbehavior =>
          pbehaviorToRequest(omit(pbehavior, ['_id'])),
        );
        const updatedPbehaviors = Object.values(this.form.changedPbehaviorsById).map(
          pbehaviorToRequest,
        );
        const removedPbehaviors = Object.values(this.form.removedPbehaviorsById);

        const [createdPbehaviorsResult, updatedPbehaviorsResult, removedPbehaviorsResult] =
          await Promise.all([
            createdPbehaviors.length && this.createPbehaviorsWithComments(createdPbehaviors),
            updatedPbehaviors.length &&
              this.updatePbehaviorsWithComments(updatedPbehaviors, originalPbehaviors),
            removedPbehaviors.length && this.removePbehaviors(removedPbehaviors),
          ]);

        await this.config.afterSubmit?.(
          createdPbehaviorsResult,
          updatedPbehaviorsResult,
          removedPbehaviorsResult,
        );

        this.closeModal();
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
