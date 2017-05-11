(() => {
    'use strict';

    const uniqueId = (() => {
        var counter = 0;

        return el => {
            if (el.prop('id')) {
                return el.prop('id');
            } else {
                return el.prop('id', `unique-id-${++counter}`)
                    .prop('id');
            }
        };
    })();

    const id = x => x;

    const addField = (container, label, placeholder) => {
        const inputNode = $('<input type="text" class="value-input"></input>')
            .prop('placeholder', placeholder);

        const labelNode = $('<label></label>')
            .text(label)
            .prop('htmlFor', uniqueId(inputNode));

        const containerNode = $('<div class="input-container"></div>');

        const removeNode = $('<input type="button" value="X">');

        removeNode.on('click', () => {
            removeNode.remove();
            labelNode.remove();
            inputNode.remove();
        });

        containerNode.append(labelNode);
        containerNode.append(inputNode);
        containerNode.append(removeNode);
        container.append(containerNode);

        return inputNode;
    };

    const error = message => {
        if (message) {
            $('#error').text(message).show();
        } else {
            $('#error').hide();
        }
    };

    $('#add-host').on('click', () => addField($('#host-names-list'), 'Domain: ', 'example.com'));
    $('#add-ip').on('click', () => addField($('#ips-list'), 'IP: ', '127.0.0.1'));

    $('#generate').on('click', () => {
        error();

        const data = {};

        data.names = $('#host-names-list .value-input')
            .toArray()
            .map(e => e.value)
            .filter(id);

        if (!data.names.length) {
            error("Enter at least one name");
            return;
        }

        data.commonName = data.names[0];

        data.ips = $('#ips-list .value-input')
            .toArray()
            .map(e => e.value)
            .filter(id);

        const input = $('<input>')
            .attr('type', 'hidden')
            .attr('name', 'data')
            .attr('value', JSON.stringify(data));

        $('<form></form>')
            .attr('action', "/create-certificate")
            .attr('method', 'POST')
            .append(input)
            .appendTo('body')
            .submit();
    });
})();
