<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline Signaler une alarme au support
      v-card-text
        v-container
          v-data-table(:items="items", :headers="tableHeaders", hide-actions)
            template(slot="items", slot-scope="{ item }")
              tr
                td {{ item.v.connector }}
                td {{ item.v.connector_name }}
                td {{ item.v.component }}
                td {{ item.v.resource }}
                td {{ item.v.creation_date | date('long', true) }}
                td {{ item.v.last_update_date | date('long', true) }}
                td {{ $t(`tables.alarmStatus.${parseInt(item.v.status.val, 10)}`) }}
                td {{ $t(`tables.alarmStates.${parseInt(item.v.state.val, 10)}`) }}
          v-select(label="Cause", :items="causes")
          v-textarea(label="Commentaire")
          v-select(label="Contacter le support", :items="mails")
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit") Notifier
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.signalAlarm,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
      tableHeaders: [
        {
          text: 'Connecteur',
          sortable: false,
        },
        {
          text: 'Nom du connecteur',
          sortable: false,
        },
        {
          text: 'Composant',
          sortable: false,
        },
        {
          text: 'Ressource',
          sortable: false,
        },
        {
          text: 'Date de création',
          sortable: false,
        },
        {
          text: 'Dernier changement',
          sortable: false,
        },
        {
          text: 'Status',
          sortable: false,
        },
        {
          text: 'State',
          sortable: false,
        },
      ],
      causes: [
        'Ne revient pas à l\'état OK',
        'Manque d\'informations',
        'Furtivité',
        'Fiche d\'aide manquante',
        'Bagottement',
      ],
      mails: [
        'example@ratp.fr', 'example2@ratp.fr',
      ],
    };
  },
};
</script>
