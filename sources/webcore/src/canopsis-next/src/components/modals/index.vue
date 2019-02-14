<template lang="pug">
  div
    modal-wrapper(
    v-for="(modal, key, index) in modals",
    :key="key",
    :modal="modal",
    :index="index",
    :dialogProps="dialogPropsMap[modal.name] || defaultDialogProps",
    )
      component(:is="modal.name", :modal="modal")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import ModalWrapper from './modal-wrapper.vue';
import CreateAckEvent from './alarm/create-ack-event.vue';
import CreateAssociateTicketEvent from './alarm/create-associate-ticket-event.vue';
import CreateCancelEvent from './alarm/create-cancel-event.vue';
import CreateChangeStateEvent from './alarm/create-change-state-event.vue';
import CreateDeclareTicketEvent from './alarm/create-declare-ticket-event.vue';
import CreateSnoozeEvent from './alarm/create-snooze-event.vue';
import VariablesHelp from './common/variables-help.vue';
import MoreInfos from './alarm/more-infos.vue';
import InfoPopupSetting from './alarm/info-popup-setting/info-popup-setting.vue';
import AddInfoPopup from './alarm/info-popup-setting/add-info-popup.vue';
import CreatePbehavior from './pbehavior/create-pbehavior.vue';
import PbehaviorList from './pbehavior/pbehavior-list.vue';
import EditLiveReporting from './alarm/edit-live-reporting.vue';
import Confirmation from './common/confirmation.vue';
import CreateWidget from './view/create-widget.vue';
import ColorPicker from './common/color-picker.vue';
import TextEditor from './common/text-editor.vue';
import TextFieldEditor from './common/text-field-editor.vue';
import CreateWatcher from './context/create-watcher.vue';
import CreateEntity from './context/create-entity.vue';
import AddEntityInfo from './context/add-entity-info.vue';
import SelectView from './view/select-view.vue';
import CreateView from './view/create-view.vue';
import CreateGroup from './view/create-group.vue';
import CreateFilter from './common/create-filter.vue';
import Watcher from './watcher/watcher.vue';
import CreateWatcherDeclareTicketEvent from './watcher/create-watcher-declare-ticket-event.vue';
import CreateWatcherPauseEvent from './watcher/create-watcher-pause-event.vue';
import ManageHistogramGroups from './stats/manage-histogram-groups.vue';
import AddStat from './stats/add-stat.vue';
import StatsDateInterval from './stats/stats-date-interval.vue';
import AlarmsList from './alarm/alarms-list.vue';
import CreateUser from './admin/create-user.vue';
import CreateRole from './admin/create-role.vue';
import CreateRight from './admin/create-right.vue';
import CreateEventFilterRule from './event-filter/create-event-filter-rule.vue';
import CreateEventFilterRulePattern from './event-filter/pattern/create-event-filter-rule-pattern.vue';
import AddEventFilterRuleToPattern from './event-filter/pattern/add-event-filter-rule-to-pattern.vue';
import EventFilterRuleActions from './event-filter/enrichment-options/event-filter-rule-actions.vue';
import EventFilterRuleExternalData from './event-filter/enrichment-options/event-filter-rule-external-data.vue';

const { mapGetters: modalMapGetters } = createNamespacedHelpers('modal');

/**
 * Wrapper for all modal windows
 */
export default {
  components: {
    ModalWrapper,
    CreateAckEvent,
    CreateAssociateTicketEvent,
    CreateCancelEvent,
    CreateChangeStateEvent,
    CreateDeclareTicketEvent,
    CreateSnoozeEvent,
    VariablesHelp,
    CreatePbehavior,
    PbehaviorList,
    EditLiveReporting,
    MoreInfos,
    InfoPopupSetting,
    AddInfoPopup,
    Confirmation,
    CreateEntity,
    CreateWatcher,
    AddEntityInfo,
    CreateWidget,
    ColorPicker,
    TextEditor,
    TextFieldEditor,
    SelectView,
    CreateView,
    CreateFilter,
    CreateGroup,
    Watcher,
    CreateWatcherDeclareTicketEvent,
    CreateWatcherPauseEvent,
    ManageHistogramGroups,
    AddStat,
    StatsDateInterval,
    AlarmsList,
    CreateUser,
    CreateRole,
    CreateRight,
    CreateEventFilterRule,
    CreateEventFilterRulePattern,
    AddEventFilterRuleToPattern,
    EventFilterRuleActions,
    EventFilterRuleExternalData,
  },
  data() {
    return {
      dialogPropsMap: {
        [MODALS.pbehaviorList]: { maxWidth: 1280, lazy: true },
        [MODALS.createWidget]: { maxWidth: 500, lazy: true },
        [MODALS.alarmsList]: { fullscreen: true, lazy: true },
        [MODALS.createFilter]: { maxWidth: 920, lazy: true },
        [MODALS.watcher]: { maxWidth: 920, lazy: true },
        [MODALS.statsDateInterval]: { maxWidth: 920, lazy: true },
      },
      defaultDialogProps: { maxWidth: 700, lazy: true },
    };
  },
  computed: {
    ...modalMapGetters(['modals']),
  },
};
</script>
