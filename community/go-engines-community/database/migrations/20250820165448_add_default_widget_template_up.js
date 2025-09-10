// Available global functions:
// genID returns a new string UUID
// isInt checks if a value is integer
// toInt transforms value to integer

if (db.widget_templates.countDocuments({title: "Default more infos", type: "alarm_more_infos"}) === 0) {
    const now = Math.ceil((new Date()).getTime() / 1000);
    db.widget_templates.insertOne({
        _id: genID(),
        title: "Default more infos",
        type: "alarm_more_infos",
        content: "{{!--\n" +
            "  TEMPLATE \"PLUS D'INFOS\"\n" +
            "  --------------------------------------------------------------------------\n" +
            "  - L'objet racine du template est `alarm`.\n" +
            "  - Les boucles `{{#each}}` utilisent `{{this.propriete}}` pour accéder aux données.\n" +
            "  - Le style est 100% en ligne et les icônes au format Material Design.\n" +
            "--}}\n" +
            "<div style=\"font-family:Arial, sans-serif;color:#2c3e50;line-height:1.5\">\n" +
            "\n" +
            "  {{!-- CARD 1: IDENTITÉ DE L'ALARME --}}\n" +
            "  <table width=\"100%\" style=\"background-color:#ffffff;border:1px solid #e0e0e0;border-radius:8px;margin-bottom:16px\">\n" +
            "    <tbody><tr>\n" +
            "      <td>\n" +
            "        <h2 style=\"margin:0 0 8px 0;font-size:1.4em;display:flex;align-items:center;gap:10px;flex-wrap:wrap\">\n" +
            "          {{state alarm.v.state.val}}\n" +
            "          <span>{{alarm.v.display_name}}</span>\n" +
            "          {{#if alarm.is_meta_alarm}}\n" +
            "            <span style=\"display:inline-flex;align-items:center;gap:6px;background-color:#e8eaf6;color:#3f51b5;border:1px solid #c5cae9;padding:4px 10px;border-radius:16px;font-size:0.9em;font-weight:500;line-height:1\">\n" +
            "              <i class=\"v-icon notranslate material-icons theme--light\" style=\"font-size:16px\">hub</i>\n" +
            "              Méta Alarme\n" +
            "            </span>\n" +
            "          {{/if}}\n" +
            "        </h2>\n" +
            "        <div style=\"font-size:0.9em\">\n" +
            "          <strong>Sévérité Max Atteinte :</strong> {{state alarm.v.max_state}}\n" +
            "          <span style=\"color:#757575;margin-left:12px\">(Initialement : {{state alarm.v.initial_state}})</span>\n" +
            "        </div>\n" +
            "      </td>\n" +
            "    </tr>\n" +
            "  </tbody></table>\n" +
            "\n" +
            "  {{!-- CARD 2: CONTEXTE ET DONNÉES CLÉS --}}\n" +
            "  <table width=\"100%\" style=\"background-color:#ffffff;border:1px solid #e0e0e0;border-radius:8px;margin-bottom:16px\">\n" +
            "    <tbody><tr>\n" +
            "      <td>\n" +
            "        <h3 style=\"margin:0 0 12px 0;font-size:1.2em\">Contexte</h3>\n" +
            "        <p style=\"margin:0 0 4px 0\"><strong>Message Actuel :</strong> {{{alarm.v.output}}}</p>\n" +
            "        <p style=\"margin:0 0 16px 0;color:#5a6d80;font-size:0.9em\"><strong>Message Initial :</strong> {{{alarm.v.initial_output}}}</p>\n" +
            "        <hr style=\"border:none;border-top:1px solid #e0e0e0;margin:16px 0\">\n" +
            "        <table width=\"100%\">\n" +
            "          <tbody><tr>\n" +
            "            <td width=\"50%\" style=\"vertical-align:top;padding-right:10px\">\n" +
            "              <ul style=\"margin:0;padding:0;list-style:none\">\n" +
            "                <li style=\"padding-bottom:4px\"><strong>Source :</strong> {{alarm.v.connector_name}}</li>\n" +
            "                <li style=\"padding-bottom:4px\"><strong>Composant :</strong> {{alarm.v.component}}</li>\n" +
            "                <li><strong>Ressource :</strong> {{alarm.v.resource}}</li>\n" +
            "              </ul>\n" +
            "            </td>\n" +
            "            <td width=\"50%\" style=\"vertical-align:top;padding-left:10px\">\n" +
            "              <ul style=\"margin:0;padding:0;list-style:none\">\n" +
            "                <li style=\"padding-bottom:4px\"><strong>Début :</strong> {{timestamp alarm.v.creation_date format=\"long\"}}</li>\n" +
            "                {{#if alarm.v.ticket}}\n" +
            "                  <li><strong>Ticket : </strong><a target=\"_blank\" href=\"{{alarm.v.ticket.ticket_url}}\">{{alarm.v.ticket.ticket}}</a></li>\n" +
            "                {{/if}}\n" +
            "              </ul>\n" +
            "            </td>\n" +
            "          </tr>\n" +
            "        </tbody></table>\n" +
            "        <div style=\"margin-top:20px\">\n" +
            "          {{#copy alarm.v.display_name}}\n" +
            "            <a style=\"text-decoration:none;display:inline-flex;align-items:center;gap:6px;background-color:#f0f0f0;border:1px solid #ddd;padding:5px 10px;border-radius:6px;cursor:pointer;font-size:0.9em;color:#2c3e50\">\n" +
            "              <i class=\"v-icon notranslate material-icons theme--light\" style=\"font-size:16px\">content_copy</i>\n" +
            "              Copier ID\n" +
            "            </a>\n" +
            "          {{/copy}}\n" +
            "          <span style=\"display:inline-block;width:8px\"></span>\n" +
            "          {{#copy (concat alarm.v.component \"/\" alarm.v.resource)}}\n" +
            "            <a style=\"text-decoration:none;display:inline-flex;align-items:center;gap:6px;background-color:#f0f0f0;border:1px solid #ddd;padding:5px 10px;border-radius:6px;cursor:pointer;font-size:0.9em;color:#2c3e50\">\n" +
            "              <i class=\"v-icon notranslate material-icons theme--light\" style=\"font-size:16px\">content_copy</i>\n" +
            "              Copier C/R\n" +
            "            </a>\n" +
            "          {{/copy}}\n" +
            "        </div>\n" +
            "      </td>\n" +
            "    </tr>\n" +
            "  </tbody></table>\n" +
            "\n" +
            "  {{!-- CARD 3: STATISTIQUES --}}\n" +
            "  <table width=\"100%\" style=\"background-color:#ffffff;border:1px solid #e0e0e0;border-radius:8px;margin-bottom:16px\">\n" +
            "    <tbody><tr>\n" +
            "      <td>\n" +
            "        <h3 style=\"margin:0 0 12px 0;font-size:1.2em\">Statistiques</h3>\n" +
            "        <table width=\"100%\" style=\"text-align:center;border-spacing:12px 0;margin-left:-12px\">\n" +
            "          <tbody><tr>\n" +
            "            <td width=\"33.3%\" style=\"background-color:#f9f9f9;padding:12px;border-radius:6px\">\n" +
            "              <div style=\"font-size:1.5em;font-weight:bold;color:#1976d2\">{{duration alarm.v.duration}}</div>\n" +
            "              <div style=\"font-size:0.9em;color:#5a6d80\">Durée Totale</div>\n" +
            "            </td>\n" +
            "            <td width=\"33.3%\" style=\"background-color:#f9f9f9;padding:12px;border-radius:6px\">\n" +
            "              <div style=\"font-size:1.5em;font-weight:bold;color:#1976d2\">{{alarm.v.events_count}}</div>\n" +
            "              <div style=\"font-size:0.9em;color:#5a6d80\">Événements</div>\n" +
            "            </td>\n" +
            "            <td width=\"33.3%\" style=\"background-color:#f9f9f9;padding:12px;border-radius:6px\">\n" +
            "              <div style=\"font-size:1.5em;font-weight:bold;color:#1976d2\">{{alarm.v.total_state_changes}}</div>\n" +
            "              <div style=\"font-size:0.9em;color:#5a6d80\">Changements d'État</div>\n" +
            "            </td>\n" +
            "          </tr>\n" +
            "        </tbody></table>\n" +
            "      </td>\n" +
            "    </tr>\n" +
            "  </tbody></table>\n" +
            "\n" +
            "  {{!-- CARD 4: SUIVI OPÉRATIONNEL --}}\n" +
            "  <table width=\"100%\" style=\"background-color:#ffffff;border:1px solid #e0e0e0;border-radius:8px;margin-bottom:16px\">\n" +
            "    <tbody><tr>\n" +
            "      <td>\n" +
            "        <h3 style=\"margin:0 0 12px 0;font-size:1.2em\">Suivi Opérationnel</h3>\n" +
            "        {{#if alarm.pbehavior}}\n" +
            "          <div style=\"background-color:#e3f2fd;color:#1565c0;padding:12px;border-radius:6px;margin-bottom:15px;display:flex;align-items:center;gap:12px\">\n" +
            "            <i class=\"v-icon notranslate material-icons theme--light\" style=\"font-size:24px;color:#1565c0\">pause_circle_outline</i>\n" +
            "            <div>\n" +
            "              <strong style=\"display:block\">Comportement périodique actif</strong>\n" +
            "              {{alarm.pbehavior.name}} ({{alarm.pbehavior.type.name}})\n" +
            "            </div>\n" +
            "          </div>\n" +
            "        {{/if}}\n" +
            "\n" +
            "        {{#if alarm.v.comments}}\n" +
            "          <h4 style=\"margin:16px 0 8px 0\">Commentaires</h4>\n" +
            "          {{#each alarm.v.comments}}\n" +
            "            <div style=\"border-bottom:1px solid #f0f0f0;padding:10px 0;display:flex;align-items:flex-start;gap:12px\">\n" +
            "              <i class=\"v-icon notranslate material-icons theme--light\" style=\"font-size:20px;color:#757575;margin-top:3px\">comment</i>\n" +
            "              <div>\n" +
            "                <strong style=\"color:#2c3e50\">{{this.a}}</strong>\n" +
            "                <span style=\"color:#5a6d80;font-size:0.85em\">(le {{timestamp this.t}})</span>\n" +
            "                <p style=\"margin:4px 0 0 0;font-style:italic\">\"{{this.m}}\"</p>\n" +
            "              </div>\n" +
            "            </div>\n" +
            "          {{/each}}\n" +
            "        {{/if}}\n" +
            "\n" +
            "        {{#if alarm.v.tickets}}\n" +
            "          <h4 style=\"margin:16px 0 8px 0\">Historique des tickets</h4>\n" +
            "          <ul style=\"margin:0;padding-left:20px;list-style-type:disc\">\n" +
            "            {{#each alarm.v.tickets}}\n" +
            "              <li style=\"padding-bottom:4px\">\n" +
            "                <a target=\"_blank\" href=\"{{this.ticket_url}}\">{{this.ticket}}</a> - créé par {{this.a}}\n" +
            "              </li>\n" +
            "            {{/each}}\n" +
            "          </ul>\n" +
            "        {{/if}}\n" +
            "      </td>\n" +
            "    </tr>\n" +
            "  </tbody></table>\n" +
            "\n" +
            "  {{!-- CARD 5: INFORMATIONS D'ENRICHISSEMENT --}}\n" +
            "  <table width=\"100%\" style=\"background-color:#ffffff;border:1px solid #e0e0e0;border-radius:8px\">\n" +
            "    <tbody><tr>\n" +
            "      <td>\n" +
            "        <h3 style=\"margin:0 0 12px 0;font-size:1.2em\">Informations d'Enrichissement</h3>\n" +
            "        <table width=\"100%\">\n" +
            "          <tbody><tr style=\"vertical-align:top\">\n" +
            "            <td width=\"33.3%\" style=\"padding-right:10px\">\n" +
            "              <strong style=\"display:block;margin-bottom:8px\">Infos Entité (Ressource)</strong>\n" +
            "              <ul style=\"margin:0;padding:0;list-style:none;font-size:0.9em\">\n" +
            "                {{#each alarm.entity.infos}}\n" +
            "                  <li style=\"padding-bottom:4px\"><strong>{{this.name}}:</strong> {{this.value}}</li>\n" +
            "                {{/each}}\n" +
            "              </ul>\n" +
            "            </td>\n" +
            "            <td width=\"33.3%\" style=\"padding-left:10px;border-left:1px solid #e0e0e0\">\n" +
            "              <strong style=\"display:block;margin-bottom:8px\">Infos Composant</strong>\n" +
            "              <ul style=\"margin:0;padding:0;list-style:none;font-size:0.9em\">\n" +
            "                {{#each alarm.entity.component_infos}}\n" +
            "                  <li style=\"padding-bottom:4px\"><strong>{{this.value}}:</strong> {{this.value}}</li>\n" +
            "                {{/each}}\n" +
            "              </ul>\n" +
            "            </td>\n" +
            "            <td width=\"33.3%\" style=\"padding-left:10px;border-left:1px solid #e0e0e0\">\n" +
            "              <strong style=\"display:block;margin-bottom:8px\">Infos Alarme (Dynamiques)</strong>\n" +
            "              {{#if alarm.v.infos}}\n" +
            "                <ul style=\"margin:0;padding:0;list-style:none;font-size:0.9em\">\n" +
            "                  {{#each alarm.v.infos}}\n" +
            "                    {{#each this}}\n" +
            "                      <li style=\"padding-bottom:4px\"><strong>{{@key}}:</strong> {{this}}</li>\n" +
            "                    {{/each}}\n" +
            "                  {{/each}}\n" +
            "                </ul>\n" +
            "              {{else}}\n" +
            "                <p style=\"margin:0;font-size:0.9em;color:#757575;font-style:italic\">Aucune info dynamique.</p>\n" +
            "              {{/if}}\n" +
            "            </td>\n" +
            "          </tr>\n" +
            "        </tbody></table>\n" +
            "      </td>\n" +
            "    </tr>\n" +
            "  </tbody></table>\n" +
            "</div>\n",
        author: "root",
        created: now,
        updated: now
    });
}
