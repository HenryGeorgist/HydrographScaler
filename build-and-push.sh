#TAG=v0.0.1
# IMAGE=lawlerseth/hydrographscalar

# docker build -t $IMAGE:$TAG .

# docker run -it --entrypoint /bin/sh $IMAGE:$TAG 

# docker push $IMAGE:$TAG

# # test
# docker run --mount type=bind,src=/home/slawler/workbench/repos/hydrographscaler/configs,dst=/workspaces/configs \
#     $IMAGE:$TAG /bin/sh -c  "./main -config=/workspaces/configs/modelpayload.json"