```console
docker network create challenge-zinc

docker run -v ~/data:/data -e ZINC_DATA_PATH="/data" -p 4080:4080 --network challenge-zinc \
    -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 \
    --name zinc public.ecr.aws/zinclabs/zinc:latest

docker run -it --rm --net challenge-zinc --name challenge-zinc-indexer javtor/challenge-zinc-indexer

docker run -it --rm --net challenge-zinc --name challenge-zinc-visualizer-back -p 3000:3000 javtor/challenge-zinc-visualizer-back

docker run -it --rm --net challenge-zinc --name challenge-zinc-visualizer-front -p 8080:8080 javtor/challenge-zinc-visualizer-front

```