package i18n

// catalog holds every message in every language. Format verbs may use
// explicit argument indexes (%[2]d) where a language orders them differently.
var catalog = map[string]map[Lang]string{
	// Requests
	"invalidBody": {
		En: "invalid request body",
		Es: "la petición no es válida",
		Eu: "eskaera ez da baliozkoa",
	},
	"needLogin": {
		En: "you need to log in",
		Es: "tienes que iniciar sesión",
		Eu: "saioa hasi behar duzu",
	},

	// Accounts
	"usernameRule": {
		En: "the username must be 3 to 20 letters, digits, '_', '.' or '-'",
		Es: "el nombre de usuario debe tener de 3 a 20 letras, dígitos, '_', '.' o '-'",
		Eu: "erabiltzaile-izenak 3 eta 20 artean letra, digitu, '_', '.' edo '-' izan behar ditu",
	},
	"emailInvalid": {
		En: "the email address is not valid",
		Es: "la dirección de correo no es válida",
		Eu: "helbide elektronikoa ez da baliozkoa",
	},
	"passwordTooShort": {
		En: "the password must have at least %d characters",
		Es: "la contraseña debe tener al menos %d caracteres",
		Eu: "pasahitzak gutxienez %d karaktere izan behar ditu",
	},
	"passwordTooLong": {
		En: "the password can't be longer than %d bytes",
		Es: "la contraseña no puede tener más de %d bytes",
		Eu: "pasahitzak ezin ditu %d byte baino gehiago izan",
	},
	"passwordsDontMatch": {
		En: "the passwords don't match",
		Es: "las contraseñas no coinciden",
		Eu: "pasahitzak ez datoz bat",
	},
	"userExists": {
		En: "a user with that email or username already exists",
		Es: "ya existe un usuario con ese correo o nombre de usuario",
		Eu: "badago helbide elektroniko edo erabiltzaile-izen hori duen erabiltzaile bat",
	},
	"createUserFailed": {
		En: "could not create the user",
		Es: "no se ha podido crear el usuario",
		Eu: "ezin izan da erabiltzailea sortu",
	},
	"credentialsRequired": {
		En: "email and password are required",
		Es: "hacen falta el correo y la contraseña",
		Eu: "helbide elektronikoa eta pasahitza behar dira",
	},
	"invalidCredentials": {
		En: "invalid email or password",
		Es: "correo o contraseña incorrectos",
		Eu: "helbide elektroniko edo pasahitz okerra",
	},
	"loginFailed": {
		En: "could not log in",
		Es: "no se ha podido iniciar sesión",
		Eu: "ezin izan da saioa hasi",
	},
	"currentPasswordWrong": {
		En: "the current password is not correct",
		Es: "la contraseña actual no es correcta",
		Eu: "uneko pasahitza ez da zuzena",
	},
	"changePasswordFailed": {
		En: "could not change the password",
		Es: "no se ha podido cambiar la contraseña",
		Eu: "ezin izan da pasahitza aldatu",
	},

	// Projects
	"loadProjectsFailed": {
		En: "could not load the projects",
		Es: "no se han podido cargar los proyectos",
		Eu: "ezin izan dira proiektuak kargatu",
	},
	"saveProjectFailed": {
		En: "could not save the project",
		Es: "no se ha podido guardar el proyecto",
		Eu: "ezin izan da proiektua gorde",
	},
	"projectNotFound": {
		En: "project not found",
		Es: "no se ha encontrado el proyecto",
		Eu: "ez da proiektua aurkitu",
	},
	"invalidProject": {
		En: "invalid project",
		Es: "el proyecto no es válido",
		Eu: "proiektua ez da baliozkoa",
	},
	"accessProjectFailed": {
		En: "could not access the project",
		Es: "no se ha podido acceder al proyecto",
		Eu: "ezin izan da proiektura sartu",
	},
	"projectNeedsName": {
		En: "the project needs a name",
		Es: "el proyecto necesita un nombre",
		Eu: "proiektuak izen bat behar du",
	},
	"projectNameTooLong": {
		En: "the project name can't be longer than %d characters",
		Es: "el nombre del proyecto no puede tener más de %d caracteres",
		Eu: "proiektuaren izenak ezin ditu %d karaktere baino gehiago izan",
	},
	"projectNeedsLayer": {
		En: "the project needs at least one layer",
		Es: "el proyecto necesita al menos una capa",
		Eu: "proiektuak geruza bat behar du gutxienez",
	},
	"tooManyLayers": {
		En: "a project can't have more than %d layers",
		Es: "un proyecto no puede tener más de %d capas",
		Eu: "proiektu batek ezin ditu %d geruza baino gehiago izan",
	},
	"layerNameTooLong": {
		En: "layer names can't be longer than %d characters",
		Es: "los nombres de las capas no pueden tener más de %d caracteres",
		Eu: "geruzen izenek ezin dituzte %d karaktere baino gehiago izan",
	},
	"layerPart": {
		En: "layer %q, part %d: %s",
		Es: "capa %q, pieza %d: %s",
		Eu: "%q geruza, %d. pieza: %s",
	},
	"tooManyParts": {
		En: "a project can't have more than %d parts",
		Es: "un proyecto no puede tener más de %d piezas",
		Eu: "proiektu batek ezin ditu %d pieza baino gehiago izan",
	},
	"unknownPartType": {
		En: "unknown part type %q",
		Es: "tipo de pieza desconocido %q",
		Eu: "pieza mota ezezaguna: %q",
	},
	"invalidColor": {
		En: "invalid color %q",
		Es: "color no válido %q",
		Eu: "kolore baliogabea: %q",
	},
	"positionRange": {
		En: "position must be between -%g and %g",
		Es: "la posición debe estar entre -%g y %g",
		Eu: "posizioak -%g eta %g artean egon behar du",
	},
	"unknownTexture": {
		En: "unknown texture %q",
		Es: "textura desconocida %q",
		Eu: "testura ezezaguna: %q",
	},
	"textureAndImage": {
		En: "a part can have a texture or an image, not both",
		Es: "una pieza puede tener una textura o una imagen, no las dos",
		Eu: "pieza batek testura bat edo irudi bat izan dezake, ez biak",
	},
	"imageLinkRule": {
		En: "image links must start with https:// and be at most %d characters",
		Es: "los enlaces de imagen deben empezar por https:// y tener como máximo %d caracteres",
		Eu: "irudi-estekek https:// honekin hasi behar dute eta gehienez %d karaktere izan",
	},
	"invalidRotation": {
		En: "invalid rotation",
		Es: "rotación no válida",
		Eu: "biraketa baliogabea",
	},
	"scaleRange": {
		En: "scale must be between %g and %g",
		Es: "el tamaño debe estar entre %g y %g",
		Eu: "tamainak %g eta %g artean egon behar du",
	},

	"invalidModelId": {
		En: "a 3D model part needs a valid model",
		Es: "una pieza de modelo 3D necesita un modelo válido",
		Eu: "3D eredu pieza batek baliozko eredu bat behar du",
	},
	"modelNoTexture": {
		En: "a 3D model keeps its own look and can't have a texture or an image",
		Es: "un modelo 3D conserva su propio aspecto y no puede tener textura ni imagen",
		Eu: "3D eredu batek bere itxura gordetzen du eta ezin du testurarik edo irudirik izan",
	},
	"modelIdOnShape": {
		En: "only 3D model parts can have a model",
		Es: "solo las piezas de modelo 3D pueden tener un modelo",
		Eu: "3D eredu piezek bakarrik izan dezakete eredu bat",
	},

	// 3D models
	"invalidModel": {
		En: "the file is not a valid .glb model",
		Es: "el archivo no es un modelo .glb válido",
		Eu: "fitxategia ez da baliozko .glb eredu bat",
	},
	"modelTooLarge": {
		En: "the model can't be larger than %d MB",
		Es: "el modelo no puede ocupar más de %d MB",
		Eu: "ereduak ezin ditu %d MB baino gehiago izan",
	},
	"modelTooManyTriangles": {
		En: "the model can't have more than %d triangles",
		Es: "el modelo no puede tener más de %d triángulos",
		Eu: "ereduak ezin ditu %d triangelu baino gehiago izan",
	},
	"modelEmpty": {
		En: "the model has no meshes",
		Es: "el modelo no tiene mallas",
		Eu: "ereduak ez du sarerik",
	},
	"modelExternalFile": {
		En: "the model uses files outside the .glb: export it with everything embedded",
		Es: "el modelo usa archivos fuera del .glb: expórtalo con todo incrustado",
		Eu: "ereduak .glb-tik kanpoko fitxategiak erabiltzen ditu: esportatu dena barruan duela",
	},
	"modelUnsupported": {
		En: "the model needs %q, which the game doesn't support: export it without compression",
		Es: "el modelo necesita %q, que el juego no admite: expórtalo sin compresión",
		Eu: "ereduak %q behar du, eta jokoak ez du onartzen: esportatu konpresiorik gabe",
	},
	"tooManyModels": {
		En: "you can't import more than %d models",
		Es: "no puedes importar más de %d modelos",
		Eu: "ezin dituzu %d eredu baino gehiago inportatu",
	},
	"saveModelFailed": {
		En: "could not save the model",
		Es: "no se ha podido guardar el modelo",
		Eu: "ezin izan da eredua gorde",
	},
	"modelNotFound": {
		En: "model not found",
		Es: "no se ha encontrado el modelo",
		Eu: "ez da eredua aurkitu",
	},

	// Game
	"shipNotAllowed": {
		En: "that ship is not available",
		Es: "esa nave no está disponible",
		Eu: "ontzi hori ez dago eskuragarri",
	},
	"shipNoVisibleParts": {
		En: "that ship has no visible parts",
		Es: "esa nave no tiene piezas visibles",
		Eu: "ontzi horrek ez du pieza ikusgairik",
	},
	"tooManyMessages": {
		En: "Disconnected: your connection sent too many messages.",
		Es: "Desconectado: tu conexión ha enviado demasiados mensajes.",
		Eu: "Deskonektatuta: zure konexioak mezu gehiegi bidali ditu.",
	},
	"guest": {
		En: "Guest",
		Es: "Invitado",
		Eu: "Gonbidatua",
	},
}
