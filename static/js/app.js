//Prompt is our JavaScript module for all alerts , notifications and custom popup dialogs.
function Prompt(){


    let toast = function(c){
        const {
            msg = "",
            icon ="success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title:msg,
            position: position,
            icon : icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,

            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function(c) {
        const {
            title = "",
            footer = "",
            msg ="",


        } = c;
        Swal.fire({
            icon : "success",
            title: title,
            footer : footer,
            text : msg,
        })
    }

    let error = function(c) {
        const {
            title = "",
            footer = "",
            msg ="",
        } = c;

        Swal.fire({
            icon : "error",
            title: title,
            footer : footer,
            text : msg,
        })
    }

    return{
        toast:toast,
        success:success,
        error:error,
        custom:custom,
    }
}

async function custom(c){
    const {
        icon = "",
        msg = "",
        title = "",
        showConfirmButton= true,
    } =c;

    const { value: result } = await Swal.fire({
        icon:icon,
        title: title,
        html: msg,
        focusConfirm: false,
        backDrop:false,
        showCancelButton : true,
        showConfirmButton: showConfirmButton,

        //This is the function that will be called when the modal is closed.
        willOpen: () => {
            if(c.willOpen !== undefined){
                c.willOpen();
            }
        },

        //This is the function that will be called when the modal is closed.
        preConfirm: () => {
            return [
                document.getElementById('start').value,
                document.getElementById('end').value,
            ]
        },
        
        //This is the function that will be called when the modal is opened.
        didOpen: () => {
            if(c.didOpen !== undefined){
                c.didOpen();
            }
        },
    })

    if(result){
        if(result.dismiss !== Swal.DismissReason.cancel){
            if(result.value!== ""){
                if(c.callback !== undefined){
                    c.callback(result);
                }
            }else{
                c.callback(false);
            }
        }else{
            c.callback(false);
        }

    }
}
