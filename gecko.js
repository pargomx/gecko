htmx.defineExtension('gecko', {
	// Usar "this" para guardar datos y funciones internamente.
	// Acceso a la "internalAPI" de htmx para usar sus utilidades.
	init: function(api) {
		// Devuelve el valor del atributo o null aplicando herencia.
		this.getClosestAttributeValue = api.getClosestAttributeValue;
	},
    onEvent: function (name, event) {
        
		// Codificar respuesta al prompt antes de enviarla en header.
        if (name === 'htmx:prompt' && event.detail.prompt) {
            this.promptEncodedVal = encodeURIComponent(event.detail.prompt);
        }
        if (name === 'htmx:configRequest'){
			if (this.promptEncodedVal) {
            	event.detail.headers['Hx-Prompt-Encoded'] = this.promptEncodedVal;
				this.promptEncodedVal = null
			}
			
			// Askfor: para solicitar al servidor una respuesta específica.
			const askforVal = this.getClosestAttributeValue(event.target,"hx-askfor")
			if (askforVal && askforVal.length > 0) {
				event.detail.headers["Hx-Askfor"] = askforVal;
			}
		}
    }
});
