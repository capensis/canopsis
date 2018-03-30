/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

 require.config({
    paths: {
        'components/component-right-checksum': 'canopsis/canopsis-rights/src/components/right-checksum/template',
        'components/component-rights-action': 'canopsis/canopsis-rights/src/components/rights-action/template',
        'components/component-rightsrenderer': 'canopsis/canopsis-rights/src/components/rightsrenderer/template',
        'components/component-rightsselector': 'canopsis/canopsis-rights/src/components/rightsselector/template',
        'editor-rights': 'canopsis/canopsis-rights/src/editors/editor-rights',
        'viewrightsform': 'canopsis/canopsis-rights/src/forms/viewrightsform/viewrightsform',
        'renderer-rights': 'canopsis/canopsis-rights/src/renderers/renderer-rights',
        'actionbutton-viewrights': 'canopsis/canopsis-rights/src/templates/actionbutton-viewrights',
        'rightschecksumbuttons': 'canopsis/canopsis-rights/src/templates/rightschecksumbuttons',
        'rightselector-itempartial': 'canopsis/canopsis-rights/src/templates/rightselector-itempartial',
        'rightselector-selecteditempartial': 'canopsis/canopsis-rights/src/templates/rightselector-selecteditempartial',

    }
});

define([
    'canopsis/canopsis-rights/src/components/right-checksum/component',
    'ehbs!components/component-right-checksum',
    'canopsis/canopsis-rights/src/components/rights-action/component',
    'ehbs!components/component-rights-action',
    'canopsis/canopsis-rights/src/components/rightsrenderer/component',
    'ehbs!components/component-rightsrenderer',
    'canopsis/canopsis-rights/src/components/rightsselector/component',
    'ehbs!components/component-rightsselector',
    'ehbs!editor-rights',
    'canopsis/canopsis-rights/src/forms/viewrightsform/controller',
    'ehbs!viewrightsform',
    'canopsis/canopsis-rights/src/objects/rightsregistry',
    'ehbs!renderer-rights',
    'canopsis/canopsis-rights/src/reopens/adapters/userview',
    'canopsis/canopsis-rights/src/reopens/controllers/application',
    'canopsis/canopsis-rights/src/reopens/mixins/crud',
    'canopsis/canopsis-rights/src/reopens/mixins/customfilterlist',
    'canopsis/canopsis-rights/src/reopens/mixins/documentation',
    'canopsis/canopsis-rights/src/reopens/mixins/showviewbutton',
    'canopsis/canopsis-rights/src/reopens/routes/application',
    'canopsis/canopsis-rights/src/reopens/routes/authenticated',
    'canopsis/canopsis-rights/src/reopens/routes/userview',
    'canopsis/canopsis-rights/src/reopens/widgets/uimaintabcollection',
    'ehbs!actionbutton-viewrights',
    'ehbs!rightschecksumbuttons',
    'ehbs!rightselector-itempartial',
    'ehbs!rightselector-selecteditempartial',
    'canopsis/canopsis-rights/src/utils/rightsflags'
], function () {
    
});
