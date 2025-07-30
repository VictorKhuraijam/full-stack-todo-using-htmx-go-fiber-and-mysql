document.addEventListener('DOMContentLoaded', function(){
    //clear form after successful submission
    document.body.addEventListener('htmx:afterSwap', function(evt) {
        if(evt.detail.target.id === 'todo-list'){
            const form = document.querySelector('.todo-form');
            form.reset();
        }
    });

    //Add loading states
    document.body.addEventListener('htmx:beforeRequest', function(evt){
        evt.target.classList.add('loading');
    })

    document.body.addEventListener('htmx:afterRequest', function(evt){
        evt.target.classList.remove('loading')
    })
    
})
